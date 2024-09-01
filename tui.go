package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav/caldav"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type tuiCalendarEvent struct {
	summary   string
	startDate time.Time
	endDate   time.Time
}

func (e tuiCalendarEvent) Title() string { return e.summary }
func (e tuiCalendarEvent) Description() string { 
		if !e.startDate.IsZero() && !e.endDate.IsZero() {
			return fmt.Sprintf("From %s to %s", e.startDate.Format(time.RFC1123), e.endDate.Format(time.RFC1123))
		} else if !e.startDate.IsZero() {
			return fmt.Sprintf("From %s", e.startDate.Format(time.RFC1123))
		} else {
			return fmt.Sprintf("Ends %s", e.endDate.Format(time.RFC1123))
		}
}
func (e tuiCalendarEvent) FilterValue() string { return e.summary }

type model struct {
    initializing bool
	error  error
	events list.Model 
}

func InitTuiModel() model {
    m := model{
        initializing: true,
        events: list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
    }
    m.events.Title = "Calendar events"

    return m
}

func (m model) Init() tea.Cmd {
	return loadEventsFromCalDAVServer
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd

	switch msg := msg.(type) {
	case ErrorMsg:
		m.error = msg.err
		return m, tea.Quit
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
    case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.events.SetSize(msg.Width-h, msg.Height-v)
	case IcalCalendarMsg:
        m.initializing = false

        var events []list.Item
		for _, calendar := range msg.calendars {
			for _, event := range calendar.Data.Events() {
				start, _ := event.DateTimeStart(time.Local)
				end, _ := event.DateTimeEnd(time.Local)
				events = append(events, tuiCalendarEvent{
					summary:   event.Props.Get(ical.PropSummary).Value,
					startDate: start,
					endDate:   end,
				})
			}
		}
        cmds = append(cmds, m.events.SetItems(events))
	}

    var cmd tea.Cmd
    m.events, cmd = m.events.Update(msg)
    cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
    if m.error != nil {
        return m.error.Error()
    }

    if m.initializing {
        return "Loading calendar events..."
    }

    return docStyle.Render(m.events.View())
}

type IcalCalendarMsg struct {
	calendars []caldav.CalendarObject
}

type ErrorMsg struct{ err error }

func (e ErrorMsg) Error() string { return e.err.Error() }

func loadEventsFromCalDAVServer() tea.Msg {
	err := godotenv.Load()
	if err != nil {
		return ErrorMsg{err: fmt.Errorf("Error loading .env file. Make sure it exists. %v", err)}
	}

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	calendarId := os.Getenv("GOOGLE_CALENDAR_ID")
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
		return ErrorMsg{err: fmt.Errorf("Could not create client: %v", err)}
	}

	principal, err := client.FindCurrentUserPrincipal(ctx)
	if err != nil {
		return ErrorMsg{err: fmt.Errorf("Could not get current user principal: %v", err)}
	}

	calHomeSet, err := client.FindCalendarHomeSet(ctx, principal)
	if err != nil {
		return ErrorMsg{err: fmt.Errorf("Could not get calendar home set: %v", err)}
	}

	calendars, err := client.FindCalendars(ctx, calHomeSet)
	if err != nil {
		return ErrorMsg{err: fmt.Errorf("Could not get calendars: %v", err)}
	}

	var queryCalendar string
	for _, calendar := range calendars {
		if calendar.Name == calendarId {
			queryCalendar = calendar.Path
		}
	}

	if queryCalendar == "" {
		return ErrorMsg{err: fmt.Errorf("Could not find calendar: %s", calendarId)}
	}

	now := time.Now()
	weekFromNow := now.AddDate(0, 0, 7)
	query := &caldav.CalendarQuery{
		CompRequest: caldav.CalendarCompRequest{
			Props: []string{"getetag", "calendar-data"},
		},
		CompFilter: caldav.CompFilter{
			Name: "VCALENDAR",
			Comps: []caldav.CompFilter{
				{
					Name:  "VEVENT",
					Start: now,
					End:   weekFromNow,
				},
			},
		},
	}
	calendarEntries, err := client.QueryCalendar(ctx, queryCalendar, query)
	if err != nil {
		return ErrorMsg{err: fmt.Errorf("Could not get calendar events: %v", err)}
	}

	return IcalCalendarMsg{
		calendars: calendarEntries,
	}
}
