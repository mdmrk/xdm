package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	Body      string
	Likes     uint32
	CreatedAt time.Time
}

const (
	BaseURL = "https://localhost:5555"
)

func main() {
	page := 0

	for {
		fmt.Println("1. Retrieve posts")
		fmt.Println("q. Exit")
		fmt.Print("> ")

		var choice string
		fmt.Scanln(&choice)

		switch choice {
		case "1":
		A:
			for {
				retrievePosts(page)

				fmt.Println("n. Next page")
				if page != 0 {
					fmt.Println("p. Previous page")
				}
				fmt.Println("b. Back to main menu")
				fmt.Print("> ")

				var pageChoice string
				fmt.Scanln(&pageChoice)

				switch pageChoice {
				case "n":
					page++
				case "p":
					if page > 0 {
						page--
					}
				case "b":
					page = 0
					break A
				default:
				}
			}
		case "q":
			return
		default:
		}
		fmt.Println()
	}
}

func retrievePosts(page int) {
	url := fmt.Sprintf("%s/posts?limit=5&offset=%d", BaseURL, page*10)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Error retrieving posts:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	var posts []Post
	err = json.Unmarshal(body, &posts)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	for _, post := range posts {
		fmt.Printf("%s\n", post.Body)
		fmt.Printf("                                %d ♥ %s\n", post.Likes, post.CreatedAt.Format(time.RFC3339))
		fmt.Println()
	}
}
