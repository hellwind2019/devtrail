package handlers

import (
	"devtrail/internal/models"
	"devtrail/internal/storage"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func HandleGitHubAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// renderTemplate(w, "github-auth.html", nil)
		code := r.URL.Query().Get("code")
		fmt.Println("Authorization code:", code)
		OAuthInfo := models.OAuthInfo{
			ClientID:     os.Getenv("GH_BASIC_CLIENT_ID"),
			ClientSecret: os.Getenv("GH_BASIC_CLIENT_SECRET"),
		}

		data := url.Values{}
		data.Set("client_id", OAuthInfo.ClientID)
		data.Set("client_secret", OAuthInfo.ClientSecret)
		data.Set("code", code)
		data.Set("redirect_uri", "http://127.0.0.1:8080/github-auth")

		req, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", strings.NewReader(data.Encode()))
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			log.Fatal(err)
		}
		token := result["access_token"].(string)
		user, sr := GetSessionUser(r, w)
		if sr {
			return
		}
		userID, _ := storage.GetUserIDByUsername(user)
		storage.SaveGitHubToken(userID, token)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

}

func MakeGitHubRequest(method, url, token string) ([]byte, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
