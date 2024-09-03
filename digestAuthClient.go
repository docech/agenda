package main

import (
	"net/http"
	"net/url"

	"github.com/emersion/go-webdav"
	"github.com/icholy/digest"
)

type DigestAuthClientConfig struct {
    Username string
    Password string
    ProxyUrl *url.URL
}

func NewDigestAuthClient(config DigestAuthClientConfig) webdav.HTTPClient {
    transport := http.Transport{}
    if config.ProxyUrl != nil {
        transport.Proxy = http.ProxyURL(config.ProxyUrl)
    }

	return &http.Client{
		Transport: &digest.Transport{
            Username: config.Username,
            Password: config.Password,
			Transport: &transport,
		},
	}
}

