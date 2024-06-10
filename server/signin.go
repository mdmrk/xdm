package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
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

func signinHandler(w http.ResponseWriter, r *http.Request) {
	var signinReq SigninRequest
	var u *User

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

	token, err := userGenerateJWT(u.Id.String())
	if err != nil {
		Error(err.Error())
		return
	}
	encode(w, r, http.StatusOK, map[string]string{"token": token, "userId": u.Id.String()})
	Info("user '%s' logged in", u.Username)
}
