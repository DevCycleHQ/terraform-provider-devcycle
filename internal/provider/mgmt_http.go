package provider

import (
	"context"
	"errors"
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
		if !shouldRetryMgmtRequest(req, resp, err) || attempt == 2 {
			return resp, err
		}

		if resp != nil && resp.Body != nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
		}

		if err := sleepWithContext(req.Context(), time.Duration(attempt+1)*500*time.Millisecond); err != nil {
			return nil, err
		}
	}

	return resp, err
}

func shouldRetryMgmtRequest(req *http.Request, resp *http.Response, err error) bool {
	switch req.Method {
	case http.MethodGet, http.MethodDelete, http.MethodHead:
	default:
		return false
	}

	if req.Context().Err() != nil {
		return false
	}

	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return false
		}
		return true
	}

	if resp == nil {
		return true
	}

	switch resp.StatusCode {
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

func sleepWithContext(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func (p *provider) setMgmtRequestHeaders(req *http.Request) {
	req.Header.Set("Authorization", p.AccessToken)
	req.Header.Set("dvc-referrer", "terraform")
	terraformVersion := p.TerraformVersion
	if terraformVersion == "" {
		terraformVersion = "unknown"
	}
	metadata := fmt.Sprintf(`{"dvc_terraform_provider_version": %q, "terraform_version": %q}`, p.version, terraformVersion)
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
