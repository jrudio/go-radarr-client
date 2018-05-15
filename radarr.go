package radarr

// Client ...
type Client struct {
	URL    string
	APIKey string
	// Timeout in seconds -- default 5
	Timeout int
}

// New creates a client to make api calls to Radarr
func New(url, apiKey string) Client {
	return Client{
		URL:     url,
		APIKey:  apiKey,
		Timeout: 5,
	}
}
