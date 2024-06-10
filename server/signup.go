package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"

	"golang.org/x/crypto/argon2"
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

	hashedPassword := hashPassword(u.Password, salt)

	_, err = db.Exec(`INSERT INTO users (alias, username, password, salt, token)
                       VALUES ($1, $2, $3, $4, $5)`,
		u.Alias, u.Username, hashedPassword, u.Salt, u.Token)
	if err != nil {
		http.Error(w, "Failed to store user", http.StatusInternalServerError)
		Error("Failed to insert user %s: %v", u.Username, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	Info("Registered user %s", u.Username)
}
