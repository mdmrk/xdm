package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
)

const (
	addr = "localhost:5556"
)

var file *os.File

func encrypt(plaintext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func handleConnection(conn net.Conn, key []byte) {
	defer conn.Close()
	buf, err := io.ReadAll(conn)
	if err != nil {
		log.Println("Failed to read from connection:", err)
		return
	}
	encryptedLog, err := encrypt(string(buf), key)
	if err != nil {
		log.Println("Failed to encrypt log:", err)
		return
	}
	if _, err := file.WriteString(encryptedLog + "\n"); err != nil {
		log.Println("Failed to write to log file:", err)
	}
}

func main() {
	key := []byte(os.Getenv("LOG_PASSWORD"))
	if len(key) == 0 {
		log.Fatalf("LOG_PASSWORD variable not defined")
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to set up listener: %v", err)
	}
	defer listener.Close()

	ex, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}

	logfileDir := filepath.Join(filepath.Dir(ex), "logs")
	err = os.MkdirAll(logfileDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	logfilePath := filepath.Join(logfileDir, fmt.Sprintf("%d.log", time.Now().UnixMilli()))
	file, err = os.OpenFile(logfilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()

	log.Printf("Logger server listening on %s", addr)
	log.Printf("Log file: %s", logfilePath)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}
		go handleConnection(conn, key)
	}
}

