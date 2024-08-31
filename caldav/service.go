package caldav

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

const (
    PROPFIND string = "PROPFIND"
    DEPTH string = "Depth"
)

type CalDAVService struct {
    serverUrl string
    client *http.Client
}

func NewCalDAVService(serverUrl string, client *http.Client) (*CalDAVService) {
    return &CalDAVService{
        serverUrl: serverUrl,
        client: client,
    }
}

func (c *CalDAVService) NewUserPrincipalRequest() (*http.Request, error) {
    body := `<?xml version="1.0" encoding="utf-8" ?>
    <d:propfind xmlns:d="DAV:">
        <d:prop>
            <d:current-user-principal />
        </d:prop>
    </d:propfind>`
    req, err := basicHttpRequest(PROPFIND, c.serverUrl, "/", body)
    if err != nil {
        return nil, err
    }
    req.Header.Set(DEPTH, "0")

    return req, nil
}

func (c *CalDAVService) NewCalendarHome() (*http.Request, error) {
    body := `<?xml version="1.0" encoding="utf-8" ?>
    <d:propfind xmlns:d="DAV:" xmlns:c="urn:ietf:params:xml:ns:caldav">
        <d:prop>
            <c:calendar-home-set />
        </d:prop>
    </d:propfind>` 
    //should be called on result href of NewUserPrincipalRequest
    req, err := basicHttpRequest(PROPFIND, c.serverUrl, "/user", body)
    if err != nil {
        return nil, err
    }
    req.Header.Set(DEPTH, "0")

    return req, nil
}

func (c *CalDAVService) NewGetAllCalendars() (*http.Request, error) {
    body := `<?xml version="1.0" encoding="utf-8" ?>
    <d:propfind xmlns:d="DAV:" xmlns:cs="http://calendarserver.org/ns/" xmlns:c="urn:ietf:params:xml:ns:caldav">
    <d:prop>
        <d:resourcetype />
        <d:displayname />
        <cs:getctag />
        <c:supported-calendar-component-set />
    </d:prop>
    </d:propfind>`
    //should be called on result href of NewCalendarHome
    req, err := basicHttpRequest(PROPFIND, c.serverUrl, "/", body)
    if err != nil {
        return nil, err
    }
    req.Header.Set(DEPTH, "1")

    return req, nil
}

func basicHttpRequest(method, baseUrl, path string, body string) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", baseUrl, path)
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, fmt.Errorf("Could not create request for url \"%s\" because: %v", url, err)
	}

	req.Header.Set("Content-Type", "application/xml; charset=utf-8")

    return req, nil
}

func (c *CalDAVService) Do(req *http.Request) ([]byte, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Request to \"%s\" failed because: %v", req.URL, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Unreadable response body: %v", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Server responded with non 2xx, 3xx status code: %d, %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

