package pirsch

import (
	"github.com/emvi/null"
	"time"
)

// Hit are the parameters to send a page hit to Pirsch.
type Hit struct {
	Hostname       string
	URL            string `json:"url"`
	IP             string `json:"ip"`
	CFConnectingIP string `json:"cf_connecting_ip"`
	XForwardedFor  string `json:"x_forwarded_for"`
	Forwarded      string `json:"forwarded"`
	XRealIP        string `json:"x_real_ip"`
	UserAgent      string `json:"user_agent"`
	AcceptLanguage string `json:"accept_language"`
	Referrer       string `json:"referrer"`
	ScreenWidth    int    `json:"screen_width"`
	ScreenHeight   int    `json:"screen_height"`
}

// Filter is used to filter statistics.
// From and To are required dates (the time is ignored).
type Filter struct {
	From                 time.Time `json:"from"`
	To                   time.Time `json:"to"`
	Path                 string    `json:"path,omitempty"`
	Pattern              string    `json:"pattern,omitempty"`
	Language             string    `json:"language,omitempty"`
	Country              string    `json:"country,omitempty"`
	Referrer             string    `json:"referrer,omitempty"`
	OS                   string    `json:"os,omitempty"`
	Browser              string    `json:"browser,omitempty"`
	Platform             string    `json:"platform,omitempty"`
	ScreenClass          string    `json:"screen_class,omitempty"`
	UTMSource            string    `json:"utm_source,omitempty"`
	UTMMedium            string    `json:"utm_medium,omitempty"`
	UTMCampaign          string    `json:"utm_campaign,omitempty"`
	UTMContent           string    `json:"utm_content,omitempty"`
	UTMTerm              string    `json:"utm_term,omitempty"`
	Limit                int       `json:"limit,omitempty"`
	IncludeAvgTimeOnPage bool      `json:"include_avg_time_on_page,omitempty"`
}

// BaseEntity contains the base data for all entities.
type BaseEntity struct {
	ID      string    `json:"id"`
	DefTime time.Time `json:"def_time"`
	ModTime time.Time `json:"mod_time"`
}

// Domain is a domain on the dashboard.
type Domain struct {
	BaseEntity

	UserID             string      `json:"user_id"`
	Hostname           string      `json:"hostname"`
	Subdomain          string      `json:"subdomain"`
	IdentificationCode string      `json:"identification_code"`
	Public             bool        `json:"public"`
	GoogleUserID       null.String `json:"google_user_id"`
	GoogleUserEmail    null.String `json:"google_user_email"`
	GSCDomain          null.String `json:"gsc_domain"`
	NewOwner           null.Int64  `json:"new_owner"`
	Timezone           null.String `json:"timezone"`
}
