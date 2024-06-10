package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"path/filepath"
	"time"
)

const (
	addr   = "localhost:5556"
	logDir = "logs/"
	aesKey = "MXi5jg4NhT1UZvtJFJHOOK3WWVHrggU="
)

var file *os.File

func encrypt(plaintext string, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return hex.EncodeToString(ciphertext), nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf, err := io.ReadAll(conn)
	if err != nil {
		log.Println("Failed to read from connection:", err)
		return
	}

	encryptedLog, err := encrypt(string(buf), aesKey)
	if err != nil {
		log.Println("Failed to encrypt log:", err)
		return
	}

	if _, err := file.WriteString(encryptedLog + "\n"); err != nil {
		log.Println("Failed to write to log file:", err)
	}
}

func main() {
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to set up listener: %v", err)
	}
	defer listener.Close()
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	var logfile_dir = path.Join(filepath.Dir(ex), "logs")
	os.Mkdir(logfile_dir, 0755)

	file, err = os.OpenFile(path.Join(logfile_dir, fmt.Sprintf("%d.txt", time.Now().UnixMilli())), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	defer file.Close()
	log.Printf("Logger server listening on %s", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}
