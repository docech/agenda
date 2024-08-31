package main

import (
	"context"
	"fmt"
	"log"

	"os"

	"github.com/emersion/go-webdav/caldav"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file. Make sure it exists")
	}

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	calDAVServerUrl := "https://apidata.googleusercontent.com/caldav/v2"

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/calendar"},
		Endpoint:     google.Endpoint,
	}

    ctx := context.Background()
    client, err := caldav.NewClient(NewOAuth2HTTPClient(config), calDAVServerUrl)
    if err != nil {
        log.Fatalf("Could not create client: %v", err)
    }

    principal, err := client.FindCurrentUserPrincipal(ctx)
    if err != nil {
        log.Fatalf("Could not get current user principal: %v", err)
    }
    fmt.Printf("Principal: %s\n", principal)

    calHomeSet, err := client.FindCalendarHomeSet(ctx, principal)
    if err != nil {
        log.Fatalf("Could not get calendar home set: %v", err)
    }
    fmt.Printf("Calendar home set: %v\n", calHomeSet)

    calendars, err := client.FindCalendars(ctx, calHomeSet)
    if err != nil {
        log.Fatalf("Could not get calendars: %v", err)
    }

    for _, calendar := range calendars {
        fmt.Printf("Calendar: %v\n", calendar)
    }

    fmt.Printf("Principal: %s\n", principal)
}
