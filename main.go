package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/docech/agenda/caldav"
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
	calendarId := os.Getenv("GOOGLE_CALENDAR_ID")
	calDAVServerUrl := fmt.Sprintf("https://apidata.googleusercontent.com/caldav/v2/%s", url.QueryEscape(calendarId))
	// calDAVServerUrl := "https://apidata.googleusercontent.com/caldav/v2"

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/calendar"},
		Endpoint:     google.Endpoint,
	}

	service := caldav.NewOAuth2CalDAVService(calDAVServerUrl, config)

// 	now := time.Now()
// 	nowPlus7Days := now.AddDate(0, 0, 7)
//
//     //Gets calendar events, REPORT, has to be sent to calendar endpoint
// 	calendarEvents := `<?xml version="1.0" encoding="utf-8" ?>
// <c:calendar-query xmlns:d="DAV:" xmlns:c="urn:ietf:params:xml:ns:caldav">
//     <d:prop>
//         <d:getetag />
//         <c:calendar-data />
//     </d:prop>
//     <c:filter>
//         <c:comp-filter name="VCALENDAR">
//             <c:comp-filter name="VEVENT">
//                 <c:time-range start="START_DATE" end="END_DATE"/>
//             </c:comp-filter>
//         </c:comp-filter>
//     </c:filter>
// </c:calendar-query>`
// 	calendarEvents = strings.ReplaceAll(calendarEvents, "START_DATE", now.Format("20060102T150405Z"))
// 	calendarEvents = strings.ReplaceAll(calendarEvents, "END_DATE", nowPlus7Days.Format("20060102T150405Z"))

    // req, err := service.NewUserPrincipalRequest()
    // req, err := service.NewCalendarHome()
    req, err := service.NewGetAllCalendars()
    if err != nil {
        log.Fatalf("Could not create request: %v", err)
    }

    fmt.Printf("Request: %v\n", req)

    response, err := service.Do(req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Response: %s\n", string(response))
}
