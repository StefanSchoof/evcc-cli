package cmd

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	evccgen "evcc-cli/internal/gen/evcc"
)

func newGeneratedClient() (*evccgen.ClientWithResponses, error) {
	baseURL := strings.TrimRight(cfg.Host, "/") + "/api"
	return evccgen.NewClientWithResponses(
		baseURL,
		evccgen.WithHTTPClient(newHTTPClient()),
	)
}

func newHTTPClient() *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = buildTLSConfig(cfg.Insecure)
	return &http.Client{
		Timeout:   cfg.Timeout,
		Transport: transport,
	}
}

func buildTLSConfig(insecure bool) *tls.Config {
	return &tls.Config{MinVersion: tls.VersionTLS12, InsecureSkipVerify: insecure}
}

func formatAPIError(status string, body []byte) error {
	return fmt.Errorf("evcc API error: %s: %s", status, strings.TrimSpace(string(body)))
}
