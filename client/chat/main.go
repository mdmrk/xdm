package main

import (
	"bufio"
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var token string
var userId string
var userAlias string
var privateKey *rsa.PrivateKey
var publicKey *rsa.PublicKey
var connectedClients = make(map[string]*rsa.PublicKey)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("1. Login")
		fmt.Println("2. Exit")
		fmt.Print("Enter your choice: ")
		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			login(scanner)
		case "2":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

func login(scanner *bufio.Scanner) {
	fmt.Print("Enter your username: ")
	scanner.Scan()
	username := scanner.Text()

	fmt.Print("Enter your password: ")
	scanner.Scan()
	password := scanner.Text()

	loginData := map[string]string{
		"username": username,
		"password": password,
	}

	jsonData, _ := json.Marshal(loginData)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Post("https://localhost:5555/signin", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Failed to login:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var loginResponse map[string]string
		json.NewDecoder(resp.Body).Decode(&loginResponse)

		token = loginResponse["token"]
		userId = loginResponse["userId"]

		fmt.Println("Login successful!")
		handleWebSocket()
	} else {
		fmt.Println("Login failed. Please try again.")
	}
}

func handleWebSocket() {
	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	conn, _, err := dialer.Dial(fmt.Sprintf("wss://localhost:5555/ws?sender=%s", userId), nil)
	if err != nil {
		fmt.Println("Failed to connect to WebSocket server:", err)
		return
	}
	defer conn.Close()

	privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Failed to generate private key:", err)
		return
	}
	publicKey = &privateKey.PublicKey

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		fmt.Println("Failed to marshal public key:", err)
		return
	}
	publicKeyBase64 := base64.StdEncoding.EncodeToString(publicKeyBytes)
	initMessage := Message{
		Type:      "publicKey",
		Sender:    userId,
		PublicKey: publicKeyBase64,
	}
	err = conn.WriteJSON(initMessage)
	if err != nil {
		fmt.Println("Failed to send public key:", err)
		return
	}

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Error reading message:", err)
				return
			}

			var msg Message
			err = json.Unmarshal(message, &msg)
			if err != nil {
				fmt.Println("Error unmarshaling message:", err)
				continue
			}

			switch msg.Type {
			case "allPublicKeys":
				connectedClients = make(map[string]*rsa.PublicKey)
				for clientID, publicKeyBase64 := range msg.AllPublicKeys {
					publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyBase64)
					if err != nil {
						fmt.Println("Error decoding public key:", err)
						continue
					}
					publicKey, err := x509.ParsePKIXPublicKey(publicKeyBytes)
					if err != nil {
						fmt.Println("Error parsing public key:", err)
						continue
					}
					connectedClients[clientID] = publicKey.(*rsa.PublicKey)
					fmt.Printf("Received and stored public key for client %s\n", clientID)
				}
			case "publicKey":
				publicKeyBytes, err := base64.StdEncoding.DecodeString(msg.PublicKey)
				if err != nil {
					fmt.Println("Error decoding public key:", err)
					continue
				}
				publicKey, err := x509.ParsePKIXPublicKey(publicKeyBytes)
				if err != nil {
					fmt.Println("Error parsing public key:", err)
					continue
				}
				connectedClients[msg.Sender] = publicKey.(*rsa.PublicKey)
				fmt.Printf("Received and stored public key for client %s\n", msg.Sender)
			case "message":
				encryptedContent, err := base64.StdEncoding.DecodeString(msg.EncryptedContent)
				if err != nil {
					fmt.Println("Error decoding encrypted content:", err)
					continue
				}
				decryptedContent, err := decryptMessage(encryptedContent, privateKey)
				if err != nil {
					fmt.Println("Error decrypting message:", err)
					continue
				}
				signature, err := base64.StdEncoding.DecodeString(msg.Signature)
				if err != nil {
					fmt.Println("Error decoding signature:", err)
					continue
				}
				publicKey, ok := connectedClients[msg.Sender]
				if !ok {
					fmt.Println("Public key not found for sender:", msg.Sender)
					continue
				}
				err = verifySignature(decryptedContent, signature, publicKey)
				if err != nil {
					fmt.Println("Error verifying signature:", err)
					continue
				}
				fmt.Printf("[%s] %s\n", msg.Sender, string(decryptedContent))
			}
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		content := scanner.Text()

		for clientID, publicKey := range connectedClients {
			if clientID == userId {
				continue
			}
			encryptedContent, err := encryptMessage([]byte(content), publicKey)
			if err != nil {
				fmt.Printf("Error encrypting message for client %s: %v\n", clientID, err)
				continue
			}
			signature, err := signMessage([]byte(content), privateKey)
			if err != nil {
				fmt.Printf("Error signing message for client %s: %v\n", clientID, err)
				continue
			}
			encryptedContentBase64 := base64.StdEncoding.EncodeToString(encryptedContent)
			signatureBase64 := base64.StdEncoding.EncodeToString(signature)

			message := Message{
				Type:             "message",
				Sender:           userId,
				Recipient:        clientID,
				EncryptedContent: encryptedContentBase64,
				Signature:        signatureBase64,
			}

			err = conn.WriteJSON(message)
			if err != nil {
				fmt.Printf("Error sending message to client %s: %v\n", clientID, err)
			}
		}
	}
}

func decryptMessage(ciphertext []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	hash := sha256.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, privateKey, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func encryptMessage(message []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	hash := sha256.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, publicKey, message, nil)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

func signMessage(message []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	hash := sha256.New()
	hash.Write(message)
	hashed := hash.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

func verifySignature(message, signature []byte, publicKey *rsa.PublicKey) error {
	hash := sha256.New()
	hash.Write(message)
	hashed := hash.Sum(nil)
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed, signature)
}

type Message struct {
	Type             string            `json:"type"`
	Sender           string            `json:"sender"`
	Recipient        string            `json:"recipient"`
	PublicKey        string            `json:"publicKey,omitempty"`
	EncryptedContent string            `json:"encryptedContent,omitempty"`
	Signature        string            `json:"signature,omitempty"`
	AllPublicKeys    map[string]string `json:"allPublicKeys,omitempty"`
}

