package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"log"
	"net/http"

	"golang.org/x/crypto/pbkdf2"
)

func signinValidator(u SigninRequest) error {
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

func signinVerifyPassword(w http.ResponseWriter, u SigninRequest) error {
	var dbu User

	err := db.QueryRow("SELECT salt, password FROM users WHERE username = $1", u.Username).Scan(&dbu.Salt, &dbu.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusUnauthorized)
		} else {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return errors.New("")
	}

	hashedPassword := hashPassword(u.Password, dbu.Salt)
	if hashedPassword != dbu.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return errors.New("Invalid credentials")
	}

	return nil
}

func decryptKey(encryptedKey []byte, password string) ([]byte, error) {
	salt := encryptedKey[:16]
	nonce := encryptedKey[16:28]
	ciphertext := encryptedKey[28:]

	derivedKey := pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)

	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}


func signinHandler(w http.ResponseWriter, r *http.Request) {
	var signinReq SigninRequest
	var u *User
	var base64EncryptedPrivateKey string
	var encryptedPrivateKey []byte
	var publicKeyPEM string

	signinReq, err := decode[SigninRequest](r)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = signinValidator(signinReq)
	if err != nil {
		http.Error(w, "Wrong parameters", http.StatusBadRequest)
		Error(err.Error())
		return
	}

	err = signinVerifyPassword(w, signinReq)
	if err != nil {
		Error(err.Error())
		return
	}

	u, err = userFindByUsername(signinReq.Username)
	if err != nil {
		Error(err.Error())
		return
	}
	err = userUpdateLastSeen(u)
	if err != nil {
		Error(err.Error())
	}

	err = db.QueryRow("SELECT private_key, public_key FROM users WHERE username = $1", signinReq.Username).Scan(&base64EncryptedPrivateKey, &publicKeyPEM)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	encryptedPrivateKey, err = base64.StdEncoding.DecodeString(base64EncryptedPrivateKey)
	privateKeyPEM, err := decryptKey(encryptedPrivateKey, signinReq.Password)
	if err != nil {
		http.Error(w, "Failed to decrypt private key", http.StatusInternalServerError)
		Error("Failed to decrypt private key: %v", err)
		return
	}

	token, err := userGenerateJWT(u.Id.String())
	if err != nil {
		Error(err.Error())
		return
	}
	encode(w, r, http.StatusOK, map[string]string{"token": token, "userId": u.Id.String(), "publicKey": publicKeyPEM, "privateKey": string(privateKeyPEM)})
	Info("user '%s' logged in", u.Username)
}
