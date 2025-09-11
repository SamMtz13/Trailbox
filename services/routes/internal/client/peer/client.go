package peer

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

// Client es un mini cliente HTTP para hablar con otro microservicio.

type Client struct {
	base   *url.URL
	client *http.Client
}

func New() (*Client, error) {
	baseStr := os.Getenv("PEER_URL")
	if baseStr == "" {
		return nil, fmt.Errorf("PEER_URL not set")
	}
	u, err := url.Parse(baseStr)
	if err != nil {
		return nil, fmt.Errorf("invalid PEER_URL: %w", err)
	}
	cli := &http.Client{
		Timeout: 8 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			IdleConnTimeout:     90 * time.Second,
			DisableCompression:  false,
			TLSHandshakeTimeout: 5 * time.Second,
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
		},
	}
	return &Client{base: u, client: cli}, nil
}

// Get hace un GET al peer en el path dado (se une a la base).
// Ejemplo: Get(ctx, "/health")
func (c *Client) Get(ctx context.Context, path string) (int, []byte, error) {
	rel, err := url.Parse(path)
	if err != nil {
		return 0, nil, fmt.Errorf("invalid path: %w", err)
	}
	u := c.base.ResolveReference(rel)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return 0, nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	return resp.StatusCode, b, err
}
