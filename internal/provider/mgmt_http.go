package provider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const mgmtAPIBaseURL = "https://api.devcycle.com"

type retryTransport struct {
	base http.RoundTripper
}

func newMgmtHTTPClient() *http.Client {
	return &http.Client{
		Transport: retryTransport{base: http.DefaultTransport},
		Timeout:   30 * time.Second,
	}
}

func (t retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	base := t.base
	if base == nil {
		base = http.DefaultTransport
	}

	var resp *http.Response
	var err error

	for attempt := 0; attempt < 3; attempt++ {
		cloned := req.Clone(req.Context())
		cloned.Header = req.Header.Clone()

		resp, err = base.RoundTrip(cloned)
		if !shouldRetryMgmtRequest(req.Method, resp, err) || attempt == 2 {
			return resp, err
		}

		if resp != nil && resp.Body != nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
		}

		time.Sleep(time.Duration(attempt+1) * 500 * time.Millisecond)
	}

	return resp, err
}

func shouldRetryMgmtRequest(method string, resp *http.Response, err error) bool {
	switch method {
	case http.MethodGet, http.MethodDelete, http.MethodHead:
	default:
		return false
	}

	if err != nil {
		return true
	}

	if resp == nil {
		return true
	}

	switch resp.StatusCode {
	case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

func (p *provider) setMgmtRequestHeaders(req *http.Request) {
	req.Header.Set("Authorization", p.AccessToken)
	req.Header.Set("dvc-referrer", "terraform")
	metadata := fmt.Sprintf(`{"dvc_terraform_provider_version": %q, "terraform_version": %q}`, p.version, "unknown")
	req.Header.Set("dvc-referrer-metadata", metadata)
	req.Header.Set("User-Agent", "terraform-provider-devcycle")
}

func (p *provider) doMgmtRequest(ctx context.Context, method, path string, query url.Values, headers map[string]string) (*http.Response, error) {
	u, err := url.Parse(mgmtAPIBaseURL + path)
	if err != nil {
		return nil, err
	}
	if len(query) > 0 {
		u.RawQuery = query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	p.setMgmtRequestHeaders(req)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := p.MgmtHTTPClient
	if client == nil {
		client = newMgmtHTTPClient()
	}

	return client.Do(req)
}
