package provider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

const mgmtAPIBaseURL = "https://api.devcycle.com"

func (p *provider) setMgmtRequestHeaders(req *http.Request) {
	req.Header.Set("Authorization", p.AccessToken)
	req.Header.Set("dvc-referrer", "terraform")
	metadata := fmt.Sprintf(`{"dvc_terraform_provider_version": %q, "terraform_version": %q}`, p.version, "unknown")
	req.Header.Set("dvc-referrer-metadata", metadata)
	req.Header.Set("User-Agent", "terraform-provider-devcycle")
}

// variablesControllerDelete removes a variable. The Management API may require
// If-Match with the latest ETag; the generated SDK delete method does not send it.
func (p *provider) variablesControllerDelete(ctx context.Context, key, projectID string, diags *diag.Diagnostics) bool {
	_, httpResp, err := p.MgmtClient.VariablesApi.VariablesControllerFindOne(ctx, key, projectID)
	if err != nil {
		diags.AddError("Client Error", fmt.Sprintf("DevCycle Terraform Error reading variable before delete: %v", err))
		return true
	}

	if httpResp.StatusCode == http.StatusNotFound {
		return false
	}

	if httpResp.StatusCode < 200 || httpResp.StatusCode > 299 {
		diags.AddError("Client Error", fmt.Sprintf("DevCycle Terraform Error: variable read returned %s", httpResp.Status))
		return true
	}

	etag := httpResp.Header.Get("ETag")
	if etag == "" {
		httpResp2, err := p.MgmtClient.VariablesApi.VariablesControllerRemove(ctx, key, projectID)
		if ret := handleDevCycleHTTP(err, httpResp2, diags); ret {
			return true
		}
		if httpResp2 != nil && httpResp2.Body != nil {
			_, _ = io.Copy(io.Discard, httpResp2.Body)
			_ = httpResp2.Body.Close()
		}
		return false
	}

	u := fmt.Sprintf("%s/v1/projects/%s/variables/%s",
		mgmtAPIBaseURL,
		url.PathEscape(projectID),
		url.PathEscape(key),
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u, nil)
	if err != nil {
		diags.AddError("Client Error", fmt.Sprintf("DevCycle Terraform Error building delete request: %v", err))
		return true
	}
	p.setMgmtRequestHeaders(req)
	req.Header.Set("If-Match", etag)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		diags.AddError("Client Error", fmt.Sprintf("DevCycle Terraform Error deleting variable: %v", err))
		return true
	}
	if resp.Body != nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}

	if resp.StatusCode == http.StatusNotFound {
		return false
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		diags.AddError("Client Error", fmt.Sprintf("DevCycle Terraform Error: %s.\nHTTP Response: %v", resp.Status, req))
		return true
	}

	return false
}
