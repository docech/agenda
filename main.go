package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type CalDAVClient struct {
	ServerUrl string
	Username  string
	Password  string
	client    *http.Client
}

func NewCalDAVClient(serverUrl string, config *oauth2.Config, token *oauth2.Token) *CalDAVClient {
	client := config.Client(context.Background(), token)

	return &CalDAVClient{
		ServerUrl: serverUrl,
		client:    client,
	}
}

func (c *CalDAVClient) Request(method, path string, body []byte) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.ServerUrl, path)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("Could not create request for url \"%s\" because: %v", url, err)
	}

	req.Header.Set("Content-Type", "application/xml; charset=utf-8")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Request to \"%s\" failed because: %v", url, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Unreadable response body: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Server responded with non 200 status code: %d, %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
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
	fmt.Printf("Please visit this URL to authorize the application: \n%s\n", authURL)

	err := openBrowser(authURL)
	if err != nil {
		fmt.Println("Unable to open browser automatically. Please open the above URL manually.")
	}

	// Wait for the code
	code := <-codeChannel

	// Shutdown the server
	ctx := context.Background()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}

	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("unable to exchange code for token: %v", err)
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
		err = fmt.Errorf("unsupported platform")
	}
	return err
}

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file. Make sure it exists")
    }

    clientID := os.Getenv("GOOGLE_CLIENT_ID")
    clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
    calendarId := os.Getenv("GOOGLE_CALENDAR_ID")
    calDAVServerUrl := fmt.Sprintf("https://apidata.googleusercontent.com/caldav/v2/%s", calendarId)

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/calendar.readonly"},
		Endpoint:     google.Endpoint,
	}

	token, err := getTokenFromWeb(config)
	if err != nil {
		log.Fatalf("Unable to get token from web: %v", err)
	}

    client := NewCalDAVClient(calDAVServerUrl, config, token)

    body := `<?xml version="1.0" encoding="utf-8" ?>
    <D:propfind xmlns:D="DAV:">
        <D:prop>
            <D:resourcetype/>
            <D:displayname/>
        </D:prop>
    </D:propfind>`

    response, err := client.Request("PROPFIND", "/events", []byte(body))
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Printf("Response: %s\n", string(response))
}
