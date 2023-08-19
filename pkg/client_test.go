package pkg

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
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
