package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"net/http"
	"sort"
	"time"

	_ "github.com/lib/pq"
)

func generateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}

// Encode keys to PEM format
func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) string {
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	return string(privateKeyPEM)
}

func encodePublicKeyToPEM(publicKey *rsa.PublicKey) (string, error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	return string(publicKeyPEM), nil
}

// Store keys in the database
func storeKeyPair(userID string, privateKeyPEM, publicKeyPEM string) error {
	_, err := db.Exec("UPDATE users SET private_key = $1, public_key = $2 WHERE id = $3", privateKeyPEM, publicKeyPEM, userID)
	return err
}

// Ensure user has keys
func ensureUserKeys(userID string) error {
	var privateKeyPEM, publicKeyPEM sql.NullString
	err := db.QueryRow("SELECT private_key, public_key FROM users WHERE id = $1", userID).Scan(&privateKeyPEM, &publicKeyPEM)
	if err != nil {
		return err
	}

	if !privateKeyPEM.Valid || !publicKeyPEM.Valid {
		privateKey, publicKey, err := generateKeyPair(2048)
		if err != nil {
			return err
		}

		privateKeyPEMStr := encodePrivateKeyToPEM(privateKey)
		publicKeyPEMStr, err := encodePublicKeyToPEM(publicKey)
		if err != nil {
			return err
		}

		err = storeKeyPair(userID, privateKeyPEMStr, publicKeyPEMStr)
		if err != nil {
			return err
		}
	}

	return nil
}

func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	senderID := r.PathValue("user_id")

	messageReq, err := decode[MessageRequest](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := ensureUserKeys(senderID); err != nil {
		http.Error(w, "Failed to ensure sender keys", http.StatusInternalServerError)
		return
	}

	if err := ensureUserKeys(messageReq.RecipientID); err != nil {
		http.Error(w, "Failed to ensure recipient keys", http.StatusInternalServerError)
		return
	}

	var publicKeyPEM string
	err = db.QueryRow("SELECT public_key FROM users WHERE id = $1", messageReq.RecipientID).Scan(&publicKeyPEM)
	if err != nil {
		http.Error(w, "Recipient not found", http.StatusNotFound)
		return
	}

	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil || block.Type != "PUBLIC KEY" {
		http.Error(w, "Invalid recipient public key", http.StatusInternalServerError)
		return
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		http.Error(w, "Invalid recipient public key", http.StatusInternalServerError)
		return
	}

	encryptedMessage, err := encryptMessage(publicKey.(*rsa.PublicKey), []byte(messageReq.Content))
	if err != nil {
		http.Error(w, "Failed to encrypt message", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO messages (sender_id, recipient_id, content) VALUES ($1, $2, $3)", senderID, messageReq.RecipientID, encryptedMessage)
	if err != nil {
		http.Error(w, "Failed to store message", http.StatusInternalServerError)
		return
	}

	encode(w, r, http.StatusOK, map[string]string{"status": "message sent"})
}

func retrieveUserMessages(w http.ResponseWriter, r *http.Request, userID string, targetUserID string) []Message {
	rows, err := db.Query(`
		SELECT sender_id, content, timestamp
		FROM messages 
		WHERE (recipient_id = $1 AND sender_id = $2) 
		ORDER BY timestamp DESC`, targetUserID, userID)
	if err != nil {
		http.Error(w, "Failed to retrieve messages", http.StatusInternalServerError)
		return []Message{}
	}
	defer rows.Close()

	var privateKeyPEM string
	err = db.QueryRow("SELECT private_key FROM users WHERE id = $1", userID).Scan(&privateKeyPEM)
	if privateKeyPEM == "" {
		return []Message{}
	}
	if err != nil && privateKeyPEM != "" {
		http.Error(w, "Failed to load private key", http.StatusInternalServerError)
		return []Message{}
	}

	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		http.Error(w, "Invalid private key", http.StatusInternalServerError)
		return []Message{}
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		http.Error(w, "Invalid private key", http.StatusInternalServerError)
		return []Message{}
	}

	var messages []Message
	for rows.Next() {
		var senderID int
		var encryptedContent []byte
		var createdAt string
		if err := rows.Scan(&senderID, &encryptedContent, &createdAt); err != nil {
			http.Error(w, "Failed to scan message", http.StatusInternalServerError)
			return []Message{}
		}

		decryptedContent, err := decryptMessage(privateKey, encryptedContent)
		if err != nil {
			http.Error(w, "Failed to decrypt message", http.StatusInternalServerError)
			return []Message{}
		}
		messages = append(messages, Message{SenderID: senderID, RecipientID: targetUserID, Content: string(decryptedContent), Timestamp: createdAt})
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Error retrieving messages", http.StatusInternalServerError)
		return []Message{}
	}

	return messages
}

func MergeAndSortMessages(m1, m2 []Message) []Message {
	merged := append(m1, m2...)
	sort.Slice(merged, func(i, j int) bool {
		t1, err1 := time.Parse(time.RFC3339, merged[i].Timestamp)
		t2, err2 := time.Parse(time.RFC3339, merged[j].Timestamp)
		if err1 != nil || err2 != nil {
			return merged[i].Timestamp < merged[j].Timestamp
		}
		return t1.Before(t2)
	})
	return merged
}

func retrieveMessagesHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	targetUserID := r.PathValue("user_id")

	m1 := retrieveUserMessages(w, r, userID, targetUserID)
	m2 := retrieveUserMessages(w, r, targetUserID, userID)

	encode(w, r, http.StatusOK, MergeAndSortMessages(m1, m2))
}

func encryptMessage(publicKey *rsa.PublicKey, message []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, message, nil)
}

func decryptMessage(privateKey *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, nil)
}
