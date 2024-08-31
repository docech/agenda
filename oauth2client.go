package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/emersion/go-webdav"
	"golang.org/x/oauth2"
)

var tokenCacheFile = "credentials.json"

func NewOAuth2HTTPClient(config *oauth2.Config) webdav.HTTPClient {
    token, err := retriveToken(tokenCacheFile, config)
    if err != nil {
        log.Fatalf("Cannot create OAuth2 HTTPClient because of: %v", err)
    }

	return config.Client(context.Background(), token)
}

func retriveToken(cacheFile string, config *oauth2.Config) (*oauth2.Token, error) {
	token, fileErr := tokenFromFile(cacheFile)

	if fileErr != nil {
		token, webErr := getTokenFromWeb(config)
		if webErr != nil {
			return nil, fmt.Errorf("Unable to get token from web nor file: %v, %v", webErr, fileErr)
		}
		saveToken(cacheFile, token)
	}

	return token, nil
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	codeChannel := make(chan string)
	randState := "super-random-str"

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != randState {
			http.Error(w, "Invalid state parameter", http.StatusBadRequest)
			return
		}

		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Code not found", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, "Authorization successful! You can close this window now.")
		codeChannel <- code
	})

	server := &http.Server{Addr: ":8080"}
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("HTTP server ListenAndServe: %v", err)
		}
	}()

	authURL := config.AuthCodeURL(randState)

    fmt.Println("Opening browser for authentication")
	err := openBrowser(authURL)
	if err != nil {
        fmt.Printf("Unable to open browser automatically. Please open URL manually: %s\n", authURL)
	}

	// Wait for the code
	code := <-codeChannel

	// Shutdown the server
	if err := server.Shutdown(context.Background()); err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}

	token, err := config.Exchange(context.TODO(), code, oauth2.AccessTypeOffline)
	if err != nil {
		return nil, fmt.Errorf("Unable to exchange code for token: %v", err)
	}

	return token, nil
}

func openBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("Cannot open browser. Unsupported platform")
	}
	return err
}

func saveToken(path string, token *oauth2.Token) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
