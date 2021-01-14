package pirsch

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	defaultBaseURL         = "https://api.pirsch.io"
	authenticationEndpoint = "/api/v1/token"
	hitEndpoint            = "/api/v1/hit"
)

var referrerQueryParams = []string{
	"ref",
	"referer",
	"referrer",
}

// Client is a client used to access the Pirsch API.
type Client struct {
	baseURL      string
	logger       *log.Logger
	clientID     string
	clientSecret string
	hostname     string
	accessToken  string
	expiresAt    time.Time
	m            sync.RWMutex
}

// ClientConfig is used to configure the Client.
type ClientConfig struct {
	// BaseURL is optional and can be used to configure a different host for the API.
	// This is usually left empty in production environments.
	BaseURL string

	// Logger is an optional logger for debugging.
	Logger *log.Logger
}

// NewClient creates a new client for given client ID, client secret, hostname, and optional configuration.
// A new client ID and secret can be generated on the Pirsch dashboard.
// The hostname must match the hostname you configured on the Pirsch dashboard (e.g. example.com).
func NewClient(clientID, clientSecret, hostname string, config *ClientConfig) *Client {
	if config == nil {
		config = &ClientConfig{
			BaseURL: defaultBaseURL,
		}
	}

	if config.BaseURL == "" {
		config.BaseURL = defaultBaseURL
	}

	return &Client{
		baseURL:      config.BaseURL,
		logger:       config.Logger,
		clientID:     clientID,
		clientSecret: clientSecret,
		hostname:     hostname,
	}
}

// Hit sends a page hit to Pirsch for given http.Request.
func (client *Client) Hit(r *http.Request) error {
	return client.performPost(client.baseURL+hitEndpoint, &Hit{
		Hostname:       client.hostname,
		URL:            r.URL.String(),
		IP:             r.RemoteAddr,
		CFConnectingIP: r.Header.Get("CF-Connecting-IP"),
		XForwardedFor:  r.Header.Get("X-Forwarded-For"),
		Forwarded:      r.Header.Get("Forwarded"),
		XRealIP:        r.Header.Get("X-Real-IP"),
		UserAgent:      r.Header.Get("User-Agent"),
		AcceptLanguage: r.Header.Get("Accept-Language"),
		Referrer:       client.getReferrerFromHeaderOrQuery(r),
	}, true)
}

func (client *Client) getReferrerFromHeaderOrQuery(r *http.Request) string {
	referrer := r.Header.Get("Referer")

	if referrer == "" {
		for _, param := range referrerQueryParams {
			referrer = r.URL.Query().Get(param)

			if referrer != "" {
				return referrer
			}
		}
	}

	return referrer
}

func (client *Client) refreshToken() error {
	body := struct {
		ClientId     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	}{
		client.clientID,
		client.clientSecret,
	}
	bodyJson, err := json.Marshal(&body)

	if err != nil {
		return err
	}

	c := http.Client{}
	resp, err := c.Post(client.baseURL+authenticationEndpoint, "application/json", bytes.NewBuffer(bodyJson))

	if err != nil {
		return err
	}

	respJson := struct {
		AccessToken string    `json:"access_token"`
		ExpiresAt   time.Time `json:"expires_at"`
	}{}

	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(&respJson); err != nil {
		return err
	}

	client.m.Lock()
	defer client.m.Unlock()
	client.accessToken = respJson.AccessToken
	client.expiresAt = respJson.ExpiresAt
	return nil
}

func (client *Client) performPost(url string, body interface{}, retry bool) error {
	reqBody, err := json.Marshal(body)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(reqBody))

	if err != nil {
		return err
	}

	client.setRequestHeaders(req)
	c := http.Client{}
	resp, err := c.Do(req)

	if err != nil {
		return err
	}

	// refresh access token and retry on 401
	if retry && resp.StatusCode == http.StatusUnauthorized {
		if err := client.refreshToken(); err != nil {
			if client.logger != nil {
				client.logger.Printf("error refreshing token: %s", err)
			}

			return err
		}

		return client.performPost(url, body, false)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return client.requestError(url, resp.StatusCode, string(body))
	}

	return nil
}

func (client *Client) setRequestHeaders(req *http.Request) {
	client.m.RLock()
	defer client.m.RUnlock()
	req.Header.Set("Authorization", "Bearer "+client.accessToken)
}

func (client *Client) requestError(url string, statusCode int, body string) error {
	if body != "" {
		return errors.New(fmt.Sprintf("%s: received status code %d on request: %s", url, statusCode, body))
	}

	return errors.New(fmt.Sprintf("%s: received status code %d on request", url, statusCode))
}
