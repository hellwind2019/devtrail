package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Commit struct {
	Sha    string `json:"sha"`
	Commit struct {
		Message string `json:"message"`
	} `json:"commit"`
}

func main123() {
	// Замінити на свій токен
	token := os.Getenv("GITHUB_TOKEN")

	owner := "hellwind2019"
	repo := "webflyx"

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", owner, repo)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "token "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var commits []Commit
	err = json.Unmarshal(body, &commits)
	if err != nil {
		fmt.Println("JSON parse error:", err)
		os.Exit(1)
	}

	for _, c := range commits {
		fmt.Printf("Commit: %s\nMessage: %s\n\n", c.Sha, c.Commit.Message)
	}
}
