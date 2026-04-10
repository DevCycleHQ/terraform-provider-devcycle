package provider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	devcyclem "github.com/devcyclehq/go-mgmt-sdk"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func (p *provider) variablesControllerDelete(ctx context.Context, key, projectID string, diags *diag.Diagnostics) bool {
	variable, httpResp, err := p.MgmtClient.VariablesApi.VariablesControllerFindOne(ctx, key, projectID)
	if httpResp != nil && httpResp.Body != nil {
		_, _ = io.Copy(io.Discard, httpResp.Body)
		_ = httpResp.Body.Close()
	}
	if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
		return false
	}
	if ret := handleDevCycleHTTP(err, httpResp, diags); ret {
		return true
	}

	if variable.Feature != "" {
		if ret := p.detachVariableFromFeature(ctx, variable, projectID, diags); ret {
			return true
		}

		variable, httpResp, err = p.waitForDetachedVariable(ctx, key, projectID)
		if httpResp != nil && httpResp.Body != nil {
			_, _ = io.Copy(io.Discard, httpResp.Body)
			_ = httpResp.Body.Close()
		}
		if ret := handleDevCycleHTTP(err, httpResp, diags); ret {
			return true
		}
		if variable.Feature != "" {
			diags.AddError("Client Error", fmt.Sprintf("DevCycle Terraform Error: variable %q is still associated with feature %q after detach", key, variable.Feature))
			return true
		}
	}

	escapedProjectID := url.PathEscape(projectID)
	escapedKey := url.PathEscape(key)
	resp, err := p.doMgmtRequest(ctx, http.MethodDelete, fmt.Sprintf("/v1/projects/%s/variables/%s", escapedProjectID, escapedKey), nil, map[string]string{
		"If-Match": "*",
	})
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
		diags.AddError("Client Error", fmt.Sprintf("DevCycle Terraform Error: %s.\nRequest URL: %s", resp.Status, responseURL(resp)))
		return true
	}

	return false
}

func (p *provider) waitForDetachedVariable(ctx context.Context, key, projectID string) (devcyclem.Variable, *http.Response, error) {
	for attempt := 0; attempt < 5; attempt++ {
		variable, httpResp, err := p.MgmtClient.VariablesApi.VariablesControllerFindOne(ctx, key, projectID)
		if err != nil || httpResp == nil || httpResp.StatusCode == http.StatusNotFound || variable.Feature == "" {
			return variable, httpResp, err
		}
		if httpResp.Body != nil {
			_, _ = io.Copy(io.Discard, httpResp.Body)
			_ = httpResp.Body.Close()
		}
		if err := sleepWithContext(ctx, time.Duration(attempt+1)*300*time.Millisecond); err != nil {
			var zero devcyclem.Variable
			return zero, nil, err
		}
	}

	return p.MgmtClient.VariablesApi.VariablesControllerFindOne(ctx, key, projectID)
}

func (p *provider) detachVariableFromFeature(ctx context.Context, variable devcyclem.Variable, projectID string, diags *diag.Diagnostics) bool {
	feature, httpResp, err := p.MgmtClient.FeaturesApi.FeaturesControllerFindOne(ctx, variable.Feature, projectID)
	if httpResp != nil && httpResp.Body != nil {
		_, _ = io.Copy(io.Discard, httpResp.Body)
		_ = httpResp.Body.Close()
	}
	if ret := handleDevCycleHTTP(err, httpResp, diags); ret {
		return true
	}

	update := devcyclem.UpdateFeatureDto{
		Name:        feature.Name,
		Key:         feature.Key,
		Description: feature.Description,
		Type_:       feature.Type_,
		Tags:        feature.Tags,
	}

	removed := false
	for _, existingVariable := range feature.Variables {
		if existingVariable.Key == variable.Key {
			removed = true
			continue
		}

		update.Variables = append(update.Variables, devcyclem.CreateVariableDto{
			Name:         existingVariable.Name,
			Description:  existingVariable.Description,
			Key:          existingVariable.Key,
			Feature:      feature.Key,
			Type_:        existingVariable.Type_,
			DefaultValue: existingVariable.DefaultValue,
		})
	}

	if !removed {
		return false
	}

	for _, variation := range feature.Variations {
		nextVariables := make(map[string]interface{}, len(variation.Variables))
		for variationKey, variationValue := range variation.Variables {
			if variationKey == variable.Key {
				continue
			}
			nextVariables[variationKey] = variationValue
		}

		update.Variations = append(update.Variations, devcyclem.FeatureVariationDto{
			Key:       variation.Key,
			Name:      variation.Name,
			Variables: nextVariables,
		})
	}

	_, httpResp, err = p.MgmtClient.FeaturesApi.FeaturesControllerUpdate(ctx, update, feature.Key, projectID)
	if httpResp != nil && httpResp.Body != nil {
		_, _ = io.Copy(io.Discard, httpResp.Body)
		_ = httpResp.Body.Close()
	}
	if ret := handleDevCycleHTTP(err, httpResp, diags); ret {
		return true
	}

	return false
}

func (p *provider) featureControllerDelete(ctx context.Context, key, projectID string, diags *diag.Diagnostics) bool {
	escapedProjectID := url.PathEscape(projectID)
	escapedKey := url.PathEscape(key)
	resp, err := p.doMgmtRequest(ctx, http.MethodDelete, fmt.Sprintf("/v1/projects/%s/features/%s", escapedProjectID, escapedKey), url.Values{
		"deleteVariables": {"true"},
	}, nil)
	if err != nil {
		diags.AddError("Client Error", fmt.Sprintf("DevCycle Terraform Error deleting feature: %v", err))
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
		diags.AddError("Client Error", fmt.Sprintf("DevCycle Terraform Error: %s.\nRequest URL: %s", resp.Status, responseURL(resp)))
		return true
	}

	return false
}

func responseURL(resp *http.Response) string {
	if resp != nil && resp.Request != nil && resp.Request.URL != nil {
		return resp.Request.URL.String()
	}
	return "<unknown>"
}
