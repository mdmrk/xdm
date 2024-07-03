package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
)

func decrypt(ciphertext string, key []byte) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(decoded) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := decoded[:aes.BlockSize]
	decoded = decoded[aes.BlockSize:]

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(decoded, decoded)

	return string(decoded), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <log_file>\n", os.Args[0])
		os.Exit(1)
	}

	logFile := os.Args[1]

	file, err := os.Open(logFile)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()

	fmt.Print("Enter the decryption key: ")
	var key string
	fmt.Scanln(&key)

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("Failed to read log file: %v", err)
		}

		line = line[:len(line)-1]

		plaintext, err := decrypt(line, []byte(key))
		if err != nil {
			log.Printf("Failed to decrypt log entry: %v", err)
			continue
		}

		fmt.Println(plaintext)
	}
}
