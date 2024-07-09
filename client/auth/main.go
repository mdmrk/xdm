package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type User struct {
	Id       string `json:"id"`
	Alias    string `json:"alias"`
	Username string `json:"username"`
	Seen     string `json:"seen"`
}

type SigninRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignupRequest struct {
	Username string `json:"username"`
	Alias    string `json:"alias"`
	Password string `json:"password"`
}

const (
	BaseURL = "https://localhost:5555"
)

var token string

func main() {
	for {
		if token == "" {
			fmt.Println("1. Register")
			fmt.Println("2. Login")
			fmt.Println("q. Exit")
		} else {
			fmt.Println("1. Register")
			fmt.Println("2. View Profile")
			fmt.Println("3. Logout")
			fmt.Println("q. Exit")
		}
		fmt.Print("> ")
		var choice string
		fmt.Scanln(&choice)
		switch choice {
		case "1":
			register()
		case "2":
			if token == "" {
				login()
			} else {
				viewProfile()
			}
		case "3":
			if token != "" {
				logout()
			} else {
				fmt.Println("Invalid choice. Please try again.")
			}
		case "q":
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
		fmt.Println()
	}
}

func register() {
	var username, alias, password string
	fmt.Print("Enter username: ")
	fmt.Scanln(&username)
	fmt.Print("Enter alias: ")
	fmt.Scanln(&alias)
	fmt.Print("Enter password: ")
	fmt.Scanln(&password)
	err := signup(username, alias, password)
	if err != nil {
		fmt.Println("Registration failed:", err)
		return
	}
	fmt.Println("Registration successful.")
}

func login() {
	var username, password string
	fmt.Print("Enter username: ")
	fmt.Scanln(&username)
	fmt.Print("Enter password: ")
	fmt.Scanln(&password)
	t, err := signin(username, password)
	if err != nil {
		fmt.Println("Login failed:", err)
		return
	}
	token = t
	fmt.Println("Login successful.")
}

func viewProfile() {
	user, err := getUser(token)
	if err != nil {
		fmt.Println("Failed to retrieve user:", err)
		return
	}
	fmt.Printf("Profile:\n")
	fmt.Printf("%s (%s)\n", user.Alias, user.Username)
	t, _ := time.Parse(time.RFC3339, user.Seen)
	fmt.Printf("Last Seen %s\n", t)
}

func logout() {
	token = ""
	fmt.Println("Logged out successfully.")
}

func signup(username, alias, password string) error {
	url := fmt.Sprintf("%s/signup", BaseURL)
	signupReq := SignupRequest{
		Username: username,
		Alias:    alias,
		Password: password,
	}
	jsonData, err := json.Marshal(signupReq)
	if err != nil {
		return err
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("registration failed with status code: %d", resp.StatusCode)
	}
	return nil
}

func signin(username, password string) (string, error) {
	url := fmt.Sprintf("%s/signin", BaseURL)
	signinReq := SigninRequest{
		Username: username,
		Password: password,
	}
	jsonData, err := json.Marshal(signinReq)
	if err != nil {
		return "", err
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var data map[string]string
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}
	token, ok := data["token"]
	if !ok {
		return "", fmt.Errorf("token not found in response")
	}
	return token, nil
}

func getUser(token string) (*User, error) {
	url := fmt.Sprintf("%s/users/me", BaseURL)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

