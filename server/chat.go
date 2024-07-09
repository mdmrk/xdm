package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	Conn       *websocket.Conn
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

var clientsMu sync.RWMutex
var connectedClients = make(map[string]*Client)

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	clients[conn] = true

	senderID := r.URL.Query().Get("sender")
	client := &Client{
		Conn:       conn,
		PublicKey:  nil,
		PrivateKey: nil,
	}
	clientsMu.Lock()
	connectedClients[senderID] = client
	clientsMu.Unlock()

	allPublicKeys := make(map[string]string)
	clientsMu.RLock()
	for clientID, c := range connectedClients {
		if clientID != senderID && c.PublicKey != nil {
			publicKeyBytes, _ := x509.MarshalPKIXPublicKey(c.PublicKey)
			allPublicKeys[clientID] = base64.StdEncoding.EncodeToString(publicKeyBytes)
		}
	}
	clientsMu.RUnlock()
	allPublicKeysMsg := Message{
		Type:          "allPublicKeys",
		AllPublicKeys: allPublicKeys,
	}
	err = conn.WriteJSON(allPublicKeysMsg)
	if err != nil {
		log.Println("Failed to send public keys:", err)
		return
	}

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading message:", err)
			delete(clients, conn)
			clientsMu.Lock()
			delete(connectedClients, senderID)
			clientsMu.Unlock()
			break
		}

		switch msg.Type {
		case "publicKey":
			publicKeyBytes, err := base64.StdEncoding.DecodeString(msg.PublicKey)
			if err != nil {
				log.Println("Failed to decode public key:", err)
				continue
			}
			publicKey, err := x509.ParsePKIXPublicKey(publicKeyBytes)
			if err != nil {
				log.Println("Failed to parse public key:", err)
				continue
			}
			clientsMu.Lock()
			connectedClients[senderID].PublicKey = publicKey.(*rsa.PublicKey)
			clientsMu.Unlock()
			log.Printf("Received and stored public key for client %s\n", senderID)

			for clientID, c := range connectedClients {
				if clientID != senderID {
					publicKeyMsg := Message{
						Type:      "publicKey",
						Sender:    senderID,
						PublicKey: msg.PublicKey,
					}
					err := c.Conn.WriteJSON(publicKeyMsg)
					if err != nil {
						log.Println("Failed to broadcast public key:", err)
					}
				}
			}

		case "message":
			clientsMu.RLock()
			recipient, ok := connectedClients[msg.Recipient]
			clientsMu.RUnlock()
			if !ok {
				log.Println("Recipient not found")
				continue
			}
			if err := recipient.Conn.WriteJSON(msg); err != nil {
				log.Println("Failed to send message:", err)
				continue
			}
		}
	}
}

func generateKeyPair(bitSize int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}

func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)
}

func encodePublicKeyToPEM(publicKey *rsa.PublicKey) ([]byte, error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: publicKeyBytes,
		},
	), nil
}

