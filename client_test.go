package pirsch

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetReferrerFromHeaderOrQuery(t *testing.T) {
	client := NewClient("", "", "", nil)
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
