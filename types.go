package pirsch

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
