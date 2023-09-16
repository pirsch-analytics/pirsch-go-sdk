package pkg

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestGetReferrerFromHeaderOrQuery(t *testing.T) {
	client := NewClient("", "", nil)
	req := httptest.NewRequest(http.MethodPost, "https://example.com/", nil)
	req.Header.Add("Referer", "header")
	assert.Equal(t, "header", client.getReferrerFromHeaderOrQuery(req))
	req = httptest.NewRequest(http.MethodPost, "https://example.com/", nil)
	assert.Empty(t, client.getReferrerFromHeaderOrQuery(req))

	for _, ref := range referrerQueryParams {
		req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("https://example.com/?%s=test", ref), nil)
		assert.Equal(t, "test", client.getReferrerFromHeaderOrQuery(req))
	}

	req = httptest.NewRequest(http.MethodPost, "https://example.com/?ref=test+space", nil)
	assert.Equal(t, "test space", client.getReferrerFromHeaderOrQuery(req))
}

func TestNewClient(t *testing.T) {
	clientID := os.Getenv("PIRSCH_CLIENT_ID")
	clientSecret := os.Getenv("PIRSCH_CLIENT_SECRET")
	baseURL := os.Getenv("PIRSCH_BASE_URL")

	if clientID != "" && clientSecret != "" {
		client := NewClient(clientID, clientSecret, &ClientConfig{
			BaseURL: baseURL,
		})
		d, err := client.Domain()
		assert.NoError(t, err)
		assert.NotNil(t, d)
	}
}

func TestGetStatsRequestURL(t *testing.T) {
	client := NewClient("", "", nil)
	url := client.getStatsRequestURL("/api/v1/test", &Filter{
		DomainID:             "o93jnhf",
		From:                 time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
		To:                   time.Date(2023, 8, 20, 0, 0, 0, 0, time.UTC),
		Start:                500,
		Scale:                ScaleDay,
		Timezone:             "Europe/Berlin",
		Path:                 "/path",
		Pattern:              "/pattern",
		EntryPath:            "/entry",
		ExitPath:             "/exit",
		Event:                "event",
		EventMetaKey:         "event_meta_key",
		EventMeta:            map[string]string{"meta": "value"},
		Language:             "en",
		Country:              "us",
		City:                 "New York",
		Referrer:             "referrer",
		ReferrerName:         "referrer_name",
		OS:                   "Windows",
		Browser:              "Firefox",
		Platform:             "desktop",
		ScreenClass:          "XXL",
		UTMSource:            "source",
		UTMMedium:            "medium",
		UTMCampaign:          "campaign",
		UTMContent:           "content",
		UTMTerm:              "term",
		CustomMetricKey:      "custom_metric_key",
		CustomMetricType:     CustomMetricTypeInteger,
		IncludeAvgTimeOnPage: true,
		Offset:               5,
		Limit:                42,
		Sort:                 "sort",
		Direction:            "asc",
		Search:               "search",
	})
	assert.Equal(t, "https://api.pirsch.io/api/v1/test?browser=Firefox&city=New+York&country=us&custom_metric_key=custom_metric_key&custom_metric_type=integer&direction=asc&entry_path=%2Fentry&event=event&event_meta_key=event_meta_key&exit_path=%2Fexit&from=2023-08-01&id=o93jnhf&include_avg_time_on_page=true&language=en&limit=42&meta_meta=value&offset=5&os=Windows&path=%2Fpath&pattern=%2Fpattern&platform=desktop&referrer=referrer&referrer_name=referrer_name&scale=day&screen_class=XXL&search=search&sort=sort&start=500&to=2023-08-20&tz=Europe%2FBerlin&utm_campaign=campaign&utm_content=content&utm_medium=medium&utm_source=source&utm_term=term", url)
}
