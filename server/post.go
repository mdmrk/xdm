package main

import (
	"net/http"
)

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	postReq, err := decode[CreatePostRequest](r)
	userID := r.Context().Value("user_id").(string)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(postReq.Body) < 1 || len(postReq.Body) > MAX_POST_LEN {
		return
	}
	_, err = db.Exec(`INSERT INTO posts (user_id,  body)
                       VALUES ($1, $2)`, userID, postReq.Body)
	if err != nil {
		http.Error(w, "Failed to store user", http.StatusInternalServerError)
		Error("Failed to store user: %v", err)
		return
	}

	encode(w, r, http.StatusOK, "Post created")
}

func likePostHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	postID := r.PathValue("post_id")

	_, err := db.Exec(`INSERT INTO likes (post_id, user_id) VALUES ($1, $2);`, postID, userID)
	if err != nil {
		http.Error(w, "Failed to like post", http.StatusInternalServerError)
		Error("Failed to like post: %v", err)
		return
	}

	encode(w, r, http.StatusOK, "Posted liked")
}

func deleteLikePostHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	postID := r.PathValue("post_id")

	q, err := db.Exec(`DELETE FROM likes WHERE post_id=$1 and user_id=$2`, postID, userID)
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

	encode(w, r, http.StatusOK, "Like removed")
error:
	http.Error(w, "Failed to remove like", http.StatusInternalServerError)
	Error("Failed to remove like: %v", err)
	return
}

func getPostHandler(w http.ResponseWriter, r *http.Request) {
	var post Post
	postID := r.PathValue("post_id")

	query := `SELECT id, user_id, body, likes, created_at FROM posts WHERE id=$1`
	err := db.QueryRow(query, postID).Scan(&post.Id, &post.UserId, &post.Body, &post.Likes, &post.CreatedAt)
	if err != nil {
		Error(err.Error())
		return
	}
	encode(w, r, http.StatusOK, post)
}

func getPostsHandler(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, user_id, body, likes, created_at FROM posts ORDER BY created_at DESC LIMIT 100`

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.Id, &p.UserId, &p.Body, &p.Likes, &p.CreatedAt); err != nil {
			http.Error(w, "Error reading rows", http.StatusInternalServerError)
			return
		}
		posts = append(posts, p)
	}
	if err = rows.Err(); err != nil {
		return
	}
	encode(w, r, http.StatusOK, posts)
}
