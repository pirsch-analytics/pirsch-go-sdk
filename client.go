package pirsch

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	defaultBaseURL          = "https://api.pirsch.io"
	authenticationEndpoint  = "/api/v1/token"
	hitEndpoint             = "/api/v1/hit"
	eventEndpoint           = "/api/v1/event"
	domainEndpoint          = "/api/v1/domain"
	sessionDurationEndpoint = "/api/v1/statistics/duration/session"
	timeOnPageEndpoint      = "/api/v1/statistics/duration/page"
	utmSourceEndpoint       = "/api/v1/statistics/utm/source"
	utmMediumEndpoint       = "/api/v1/statistics/utm/medium"
	utmCampaignEndpoint     = "/api/v1/statistics/utm/campaign"
	utmContentEndpoint      = "/api/v1/statistics/utm/content"
	utmTermEndpoint         = "/api/v1/statistics/utm/term"
	visitorsEndpoint        = "/api/v1/statistics/visitor"
	pagesEndpoint           = "/api/v1/statistics/page"
	conversionGoalsEndpoint = "/api/v1/statistics/goals"
	eventsEndpoint          = "/api/v1/statistics/events"
	eventMetadataEndpoint   = "/api/v1/statistics/event/meta"
	growthRateEndpoint      = "/api/v1/statistics/growth"
	activeVisitorsEndpoint  = "/api/v1/statistics/active"
	timeOfDayEndpoint       = "/api/v1/statistics/hours"
	languageEndpoint        = "/api/v1/statistics/language"
	referrerEndpoint        = "/api/v1/statistics/referrer"
	osEndpoint              = "/api/v1/statistics/os"
	browserEndpoint         = "/api/v1/statistics/browser"
	countryEndpoint         = "/api/v1/statistics/country"
	platformEndpoint        = "/api/v1/statistics/platform"
	screenEndpoint          = "/api/v1/statistics/screen"
	keywordsEndpoint        = "/api/v1/statistics/keywords"
	requestRetries          = 5
)

var referrerQueryParams = []string{
	"ref",
	"referer",
	"referrer",
	"source",
	"utm_source",
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

// HitOptions optional parameters to send with the hit request.
type HitOptions struct {
	ScreenWidth  int
	ScreenHeight int
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
	return client.HitWithOptions(r, nil)
}

// HitWithOptions sends a page hit to Pirsch for given http.Request and options.
func (client *Client) HitWithOptions(r *http.Request, options *HitOptions) error {
	if r.Header.Get("DNT") == "1" {
		return nil
	}

	if options == nil {
		options = new(HitOptions)
	}

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
		ScreenWidth:    options.ScreenWidth,
		ScreenHeight:   options.ScreenHeight,
	}, requestRetries)
}

// Event sends an event to Pirsch for given http.Request.
func (client *Client) Event(name string, durationSeconds int, meta map[string]string, r *http.Request) error {
	return client.EventWithOptions(name, durationSeconds, meta, r, nil)
}

// EventWithOptions sends an event to Pirsch for given http.Request and options.
func (client *Client) EventWithOptions(name string, durationSeconds int, meta map[string]string, r *http.Request, options *HitOptions) error {
	if r.Header.Get("DNT") == "1" {
		return nil
	}

	if options == nil {
		options = new(HitOptions)
	}

	metaKeys := make([]string, 0)
	metaValues := make([]string, 0)

	for k, v := range meta {
		metaKeys = append(metaKeys, k)
		metaValues = append(metaKeys, v)
	}

	return client.performPost(client.baseURL+eventEndpoint, &Event{
		Name:            name,
		DurationSeconds: durationSeconds,
		MetaKeys:        metaKeys,
		MetaValues:      metaValues,
		Hit: Hit{
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
			ScreenWidth:    options.ScreenWidth,
			ScreenHeight:   options.ScreenHeight,
		},
	}, requestRetries)
}

// Domain returns the domain for this client.
func (client *Client) Domain() (*Domain, error) {
	domains := make([]Domain, 0, 1)

	if err := client.performGet(client.baseURL+domainEndpoint, nil, requestRetries, &domains); err != nil {
		return nil, err
	}

	if len(domains) != 1 {
		return nil, errors.New("domain not found")
	}

	return &domains[0], nil
}

// SessionDuration returns the session duration grouped by day.
func (client *Client) SessionDuration(filter *Filter) ([]TimeSpentStats, error) {
	stats := make([]TimeSpentStats, 0)

	if err := client.performGet(client.getStatsRequestURL(sessionDurationEndpoint, filter.DomainID), filter, requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// TimeOnPage returns the time spent on pages.
func (client *Client) TimeOnPage(filter *Filter) ([]TimeSpentStats, error) {
	stats := make([]TimeSpentStats, 0)

	if err := client.performGet(client.getStatsRequestURL(timeOnPageEndpoint, filter.DomainID), filter, requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// UTMSource returns the utm sources.
func (client *Client) UTMSource(filter *Filter) ([]UTMSourceStats, error) {
	stats := make([]UTMSourceStats, 0)

	if err := client.performGet(client.getStatsRequestURL(utmSourceEndpoint, filter.DomainID), filter, requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// UTMMedium returns the utm medium.
func (client *Client) UTMMedium(filter *Filter) ([]UTMMediumStats, error) {
	stats := make([]UTMMediumStats, 0)

	if err := client.performGet(client.getStatsRequestURL(utmMediumEndpoint, filter.DomainID), filter, requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// UTMCampaign returnst he utm campaigns.
func (client *Client) UTMCampaign(filter *Filter) ([]UTMCampaignStats, error) {
	stats := make([]UTMCampaignStats, 0)

	if err := client.performGet(client.getStatsRequestURL(utmCampaignEndpoint, filter.DomainID), filter, requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// UTMContent returns the utm content.
func (client *Client) UTMContent(filter *Filter) ([]UTMContentStats, error) {
	stats := make([]UTMContentStats, 0)

	if err := client.performGet(client.getStatsRequestURL(utmContentEndpoint, filter.DomainID), filter, requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// UTMTerm returns the utm term.
func (client *Client) UTMTerm(filter *Filter) ([]UTMTermStats, error) {
	stats := make([]UTMTermStats, 0)

	if err := client.performGet(client.getStatsRequestURL(utmTermEndpoint, filter.DomainID), filter, requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// Visitors returns the visitor statistics grouped by day.
func (client *Client) Visitors(filter *Filter) ([]VisitorStats, error) {
	stats := make([]VisitorStats, 0)

	if err := client.performGet(client.getStatsRequestURL(visitorsEndpoint, filter.DomainID), filter, requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// Pages returns the page statistics grouped by page.
func (client *Client) Pages(filter *Filter) ([]PageStats, error) {
	stats := make([]PageStats, 0)

	if err := client.performGet(client.getStatsRequestURL(pagesEndpoint, filter.DomainID), filter, requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// ConversionGoals returns all conversion goals.
func (client *Client) ConversionGoals(filter *Filter) ([]ConversionGoal, error) {
	stats := make([]ConversionGoal, 0)

	if err := client.performGet(client.getStatsRequestURL(conversionGoalsEndpoint, filter.DomainID), filter, requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// Events returns all events.
func (client *Client) Events(filter *Filter) ([]EventStats, error) {
	stats := make([]EventStats, 0)

	if err := client.performGet(client.getStatsRequestURL(eventsEndpoint, filter.DomainID), filter, requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// EventMetadata returns the metadata values for an event and key.
func (client *Client) EventMetadata(filter *Filter) ([]EventStats, error) {
	stats := make([]EventStats, 0)

	if err := client.performGet(client.getStatsRequestURL(eventMetadataEndpoint, filter.DomainID), filter, requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// Growth returns the growth rates for visitors, bounces, ...
func (client *Client) Growth(filter *Filter) (*Growth, error) {
	growth := new(Growth)

	if err := client.performGet(client.getStatsRequestURL(growthRateEndpoint, filter.DomainID), filter, requestRetries, growth); err != nil {
		return nil, err
	}

	return growth, nil
}

// ActiveVisitors returns the active visitors and what pages they're on.
func (client *Client) ActiveVisitors(filter *Filter) (*ActiveVisitorsData, error) {
	active := new(ActiveVisitorsData)

	if err := client.performGet(client.getStatsRequestURL(activeVisitorsEndpoint, filter.DomainID), filter, requestRetries, active); err != nil {
		return nil, err
	}

	return active, nil
}

// TimeOfDay returns the number of unique visitors grouped by time of day.
func (client *Client) TimeOfDay(filter *Filter) ([]VisitorHourStats, error) {
	stats := make([]VisitorHourStats, 0)

	if err := client.performGet(client.getStatsRequestURL(timeOfDayEndpoint, filter.DomainID), filter, requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// Languages returns language statistics.
func (client *Client) Languages(filter *Filter) ([]LanguageStats, error) {
	stats := make([]LanguageStats, 0)

	if err := client.performGet(client.getStatsRequestURL(languageEndpoint, filter.DomainID), filter, requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// Referrer returns referrer statistics.
func (client *Client) Referrer(filter *Filter) ([]ReferrerStats, error) {
	stats := make([]ReferrerStats, 0)

	if err := client.performGet(client.getStatsRequestURL(referrerEndpoint, filter.DomainID), filter, requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// OS returns operating system statistics.
func (client *Client) OS(filter *Filter) ([]OSStats, error) {
	stats := make([]OSStats, 0)

	if err := client.performGet(client.getStatsRequestURL(osEndpoint, filter.DomainID), filter, requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// Browser returns browser statistics.
func (client *Client) Browser(filter *Filter) ([]BrowserStats, error) {
	stats := make([]BrowserStats, 0)

	if err := client.performGet(client.getStatsRequestURL(browserEndpoint, filter.DomainID), filter, requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// Country returns country statistics.
func (client *Client) Country(filter *Filter) ([]CountryStats, error) {
	stats := make([]CountryStats, 0)

	if err := client.performGet(client.getStatsRequestURL(countryEndpoint, filter.DomainID), filter, requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// Platform returns the platforms used by visitors.
func (client *Client) Platform(filter *Filter) (*PlatformStats, error) {
	platforms := new(PlatformStats)

	if err := client.performGet(client.getStatsRequestURL(platformEndpoint, filter.DomainID), filter, requestRetries, platforms); err != nil {
		return nil, err
	}

	return platforms, nil
}

// Screen returns the screen classes used by visitors.
func (client *Client) Screen(filter *Filter) ([]ScreenClassStats, error) {
	stats := make([]ScreenClassStats, 0)

	if err := client.performGet(client.getStatsRequestURL(screenEndpoint, filter.DomainID), filter, requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// Keywords returns the Google keywords, rank, and CTR.
func (client *Client) Keywords(filter *Filter) ([]Keyword, error) {
	stats := make([]Keyword, 0)

	if err := client.performGet(client.getStatsRequestURL(keywordsEndpoint, filter.DomainID), filter, requestRetries, &stats); err != nil {
		return nil, err
	}

	return stats, nil
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
	client.m.Lock()
	defer client.m.Unlock()

	// check token has expired or is about to expire soon
	if client.expiresAt.After(time.Now().UTC().Add(-time.Minute)) {
		return nil
	}

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

	client.accessToken = respJson.AccessToken
	client.expiresAt = respJson.ExpiresAt
	return nil
}

func (client *Client) performPost(url string, body interface{}, retry int) error {
	reqBody, err := json.Marshal(body)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(reqBody))

	if err != nil {
		return err
	}

	client.m.RLock()
	req.Header.Set("Authorization", "Bearer "+client.accessToken)
	client.m.RUnlock()
	c := http.Client{}
	resp, err := c.Do(req)

	if err != nil {
		return err
	}

	// refresh access token and retry on 401
	if retry > 0 && resp.StatusCode == http.StatusUnauthorized {
		time.Sleep(time.Millisecond * time.Duration((requestRetries-retry)*100+50))

		if err := client.refreshToken(); err != nil {
			if client.logger != nil {
				client.logger.Printf("error refreshing token: %s", err)
			}

			return err
		}

		return client.performPost(url, body, retry-1)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return client.requestError(url, resp.StatusCode, string(body))
	}

	return nil
}

func (client *Client) performGet(url string, body interface{}, retry int, result interface{}) error {
	if retry > 0 && client.accessToken == "" {
		time.Sleep(time.Millisecond * time.Duration((requestRetries-retry)*100+50))

		if err := client.refreshToken(); err != nil {
			if client.logger != nil {
				client.logger.Printf("error refreshing token: %s", err)
			}

			return err
		}

		return client.performGet(url, body, retry-1, result)
	}

	reqBody, err := json.Marshal(body)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", url, bytes.NewReader(reqBody))

	if err != nil {
		return err
	}

	client.m.RLock()
	req.Header.Set("Authorization", "Bearer "+client.accessToken)
	req.Header.Set("Content-Type", "application/json")
	client.m.RUnlock()
	c := http.Client{}
	resp, err := c.Do(req)

	if err != nil {
		return err
	}

	// refresh access token and retry on 401
	if retry > 0 && resp.StatusCode == http.StatusUnauthorized {
		time.Sleep(time.Millisecond * time.Duration((requestRetries-retry)*100+50))

		if err := client.refreshToken(); err != nil {
			if client.logger != nil {
				client.logger.Printf("error refreshing token: %s", err)
			}

			return err
		}

		return client.performGet(url, body, retry-1, result)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return client.requestError(url, resp.StatusCode, string(body))
	}

	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(result); err != nil {
		return err
	}

	return nil
}

func (client *Client) requestError(url string, statusCode int, body string) error {
	if body != "" {
		return errors.New(fmt.Sprintf("%s: received status code %d on request: %s", url, statusCode, body))
	}

	return errors.New(fmt.Sprintf("%s: received status code %d on request", url, statusCode))
}

func (client *Client) getStatsRequestURL(endpoint, id string) string {
	return fmt.Sprintf("%s%s?id=%s", client.baseURL, endpoint, id)
}
