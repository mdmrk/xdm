package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func userUpdateLastSeen(u *User) error {
	_, err := db.Exec("UPDATE users SET seen = NOW() WHERE username = $1", u.Username)
	return err
}

func userGenerateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("your-256-bit-secret"))
	return tokenString, err
}

func userFindByUsername(username string) (*User, error) {
	var user User

	query := `SELECT id, username, password, alias FROM users WHERE username = $1`
	err := db.QueryRow(query, username).Scan(&user.Id, &user.Username, &user.Password, &user.Alias)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("User not found")
		}
		return nil, err
	}

	return &user, nil
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	userID := r.PathValue("user_id")

	query := `SELECT id, alias, username, seen FROM users WHERE id=$1`
	err := db.QueryRow(query, userID).Scan(&user.Id, &user.Alias, &user.Username, &user.Seen)
	if err != nil {
		Error(err.Error())
		return
	}
	encode(w, r, http.StatusOK, user)
}

func userFollowHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	followedID := r.PathValue("user_id")

	query := `INSERT INTO follows (follower_id, followed_id) VALUES ($1, $2)`
	_, err := db.Exec(query, userID, followedID)
	if err != nil {
		http.Error(w, "Failed to follow user", http.StatusInternalServerError)
		Error("Failed to follow user")
		return
	}
	encode(w, r, http.StatusOK, "User followed")
}

func userUnfollowHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	followedID := r.PathValue("user_id")

	q, err := db.Exec(`DELETE FROM follows WHERE follower_id=$1 and followed_id=$2`, userID, followedID)
	if err != nil {
		goto error
	} else {
		count, err := q.RowsAffected()
		if err == nil {
			if count == 0 {
				goto error
			}
		} else {
			goto error
		}
	}

	encode(w, r, http.StatusOK, "User unfollowed")
error:
	http.Error(w, "Failed to unfollow", http.StatusInternalServerError)
	Error("Failed to unfollow: %v", err)
	return
}

func userLikesHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")
	query := `SELECT post_id FROM likes WHERE user_id=$1`

	rows, err := db.Query(query, userID)
	if err != nil {
		http.Error(w, "Query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var postIDs []string
	for rows.Next() {
		var postID string
		if err := rows.Scan(&postID); err != nil {
			http.Error(w, "Error reading rows", http.StatusInternalServerError)
			return
		}
		postIDs = append(postIDs, postID)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, "Error after reading rows", http.StatusInternalServerError)
		return
	}

	encode(w, r, http.StatusOK, postIDs)
}
