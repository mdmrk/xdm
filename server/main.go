package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
)

var db *sql.DB

func encode[T any](w http.ResponseWriter, _ *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func decode[T any](r *http.Request) (T, error) {
	var v T
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

func extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractToken(r)
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		keyFunc := func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(JWT_SECRET_KEY), nil
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		next(w, r.WithContext(ctx))
	}
}

func addRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /signup", signupHandler)
	mux.HandleFunc("POST /signin", signinHandler)
	mux.HandleFunc("POST /users/{user_id}/follow", authMiddleware(userFollowHandler))
	mux.HandleFunc("DELETE /users/{user_id}/follow", authMiddleware(userUnfollowHandler))
	mux.HandleFunc("POST /posts", authMiddleware(createPostHandler))
	mux.HandleFunc("GET /posts", getPostsHandler)
	mux.HandleFunc("GET /users/{user_id}", getUserHandler)
	mux.HandleFunc("GET /posts/{post_id}", getPostHandler)
	mux.HandleFunc("POST /posts/{post_id}/like", authMiddleware(likePostHandler))
	mux.HandleFunc("DELETE /posts/{post_id}/like", authMiddleware(deleteLikePostHandler))
	mux.HandleFunc("GET /users/{user_id}/like", authMiddleware(userLikesHandler))
	mux.HandleFunc("POST /users/{user_id}/chat", authMiddleware(sendMessageHandler))
	mux.HandleFunc("GET /users/{user_id}/chat", authMiddleware(retrieveMessagesHandler))
}

func connectDB() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", DBHOST, DBPORT, DBUSER, DBPASS, DBNAME)
	Info("connecting to database %s", DBHOST)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	Info("testing connection")
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	Info("database connected")
	return db, nil
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func run() error {
	var err error
	addr := fmt.Sprintf(":%d", SERVER_PORT)

	db, err = connectDB()
	if err != nil {
		Error(err.Error())
		Error("Unsuccessful database connection attempt")
		os.Exit(1)
	}
	defer db.Close()

	mux := http.NewServeMux()
	addRoutes(mux)
	handler := corsMiddleware(mux)

	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS13,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		InsecureSkipVerify:       true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	Info("listening on addr %s", addr)
	Fatal(srv.ListenAndServeTLS("crypto/server.crt", "crypto/server.key").Error())
	return nil
}

func main() {
	if err := run(); err != nil {
		Fatal(err.Error())
	}
}
