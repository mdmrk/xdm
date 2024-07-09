package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io"
	"net/http"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/pbkdf2"
)

func generateSalt(size int) ([]byte, error) {
	salt := make([]byte, size)
	_, err := rand.Read(salt)

	if err != nil {
		return nil, err
	}
	return salt, nil
}

func hashPassword(password string, salt []byte) string {
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return base64.StdEncoding.EncodeToString(hash)
}

func signupValidator(u User) error {
	if len(u.Alias) > MAX_ALIAS_LEN || len(u.Alias) < MIN_ALIAS_LEN {
		return errors.New("Alias has wrong length")
	}
	if !VALID_CHARS.MatchString(u.Alias) {
		return errors.New("Alias contains invalid characters; only letters, numbers, '.', '_', and '-' are allowed")
	}

	if len(u.Username) > MAX_USERNAME_LEN || len(u.Username) < MIN_USERNAME_LEN {
		return errors.New("Username has wrong length")
	}
	if !VALID_CHARS.MatchString(u.Username) {
		return errors.New("Username contains invalid characters; only letters, numbers, '.', '_', and '-' are allowed")
	}

	if len(u.Password) > MAX_PASSWORD_LEN || len(u.Password) < MIN_PASSWORD_LEN {
		return errors.New("Password has wrong length")
	}
	if !VALID_CHARS.MatchString(u.Password) {
		return errors.New("Password contains invalid characters; only letters, numbers, '.', '_', and '-' are allowed")
	}

	return nil
}

func encryptKey(key []byte, password string, salt []byte) ([]byte, error) {
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	derivedKey := pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)

	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	ciphertext := aesgcm.Seal(nil, nonce, key, nil)
	return append(salt, append(nonce, ciphertext...)...), nil
}


func signupHandler(w http.ResponseWriter, r *http.Request) {
	var u User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = signupValidator(u)
	if err != nil {
		http.Error(w, "Wrong parameters", http.StatusBadRequest)
		return
	}

	salt, err := generateSalt(16)
	if err != nil {
		http.Error(w, "Failed to generate salt", http.StatusInternalServerError)
		return
	}

	u.Salt = salt

	// Generate a new key pair
	privateKey, publicKey, err := generateKeyPair(2048)
	if err != nil {
		http.Error(w, "Failed to generate key pair", http.StatusInternalServerError)
		Error("Failed to generate key pair: %v", err)
		return
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(publicKey),
	})

	encryptedPrivateKey, err := encryptKey(privateKeyPEM, u.Password, u.Salt)
	if err != nil {
		http.Error(w, "Failed to encrypt private key", http.StatusInternalServerError)
		Error("Failed to encrypt private key: %v", err)
		return
	}
	base64EncryptedPrivateKey := base64.StdEncoding.EncodeToString(encryptedPrivateKey)
	hashedPassword := hashPassword(u.Password, salt)

	_, err = db.Exec(`INSERT INTO users (alias, username, password, salt, private_key, public_key)
                       VALUES ($1, $2, $3, $4, $5, $6)`,
		u.Alias, u.Username, hashedPassword, u.Salt, base64EncryptedPrivateKey, publicKeyPEM)
	if err != nil {
		http.Error(w, "Failed to insert user", http.StatusInternalServerError)
		Error("Failed to insert user %s: %v", u.Username, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	Info("Registered user %s", u.Username)
}
