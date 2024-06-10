package main

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

type User struct {
	Id       uuid.UUID
	Alias    string
	Username string
	Password string
	Salt     []byte
	Token    []byte
	Seen     time.Time
}

type Post struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	Body      string
	Likes     uint32
	CreatedAt time.Time
}

func (p Post) FilterValue() string { return p.Body }

type SigninRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignupRequest struct {
	Username string `json:"username"`
	Alias    string `json:"alias"`
	Password string `json:"password"`
}

type CreatePostRequest struct {
	Body string `json:"body"`
}

type MessageRequest struct {
	RecipientID string `json:"recipient_id"`
	Content     string `json:"content"`
}

type Message struct {
	RecipientID string `json:"recipient_id"`
	SenderID    int    `json:"sender_id"`
	Content     string `json:"content"`
	Timestamp   string `json:"timestamp"`
}
