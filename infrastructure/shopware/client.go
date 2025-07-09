package shopware

import (
	"context"
	"net/http"
	"strings"

	"golang.org/x/oauth2/clientcredentials"
)

const (
    TokenEndpoint = "/api/oauth/token"
)

type Client struct {
	HttpClient *http.Client
	BaseURL string
}

func NewSwClientFromIntegration(ctx context.Context, clientID string, clientSecret string, baseURL string) *Client {
	baseURL = strings.TrimRight(baseURL, "/")
	tokenURL := baseURL + TokenEndpoint
	cfg := clientcredentials.Config{
		ClientID:    clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
	}

	tcn, err := cfg.Token(ctx)
	if tcn == nil && err != nil {
		panic("failed to obtain token from client credentials config")
	}

	oAusClient := cfg.Client(ctx)
	
	return &Client{
		HttpClient: oAusClient,
		BaseURL: baseURL,
	}
}
