package alibabacloud

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	APIMetaBaseURL = "https://api.aliyun.com/meta/v1"
)

// APIMetaClient provides access to AlibabaCloud API metadata
type APIMetaClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewAPIMetaClient creates a new API meta client
func NewAPIMetaClient() *APIMetaClient {
	return &APIMetaClient{
		baseURL: APIMetaBaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ProductInfo represents product information
type ProductInfo struct {
	Code           string `json:"code"`
	Name           string `json:"name"`
	DefaultVersion string `json:"defaultVersion"`
	Style          string `json:"style"`
}

// APIInfo represents API information
type APIInfo struct {
	Name       string                 `json:"name"`
	Summary    string                 `json:"summary"`
	Path       string                 `json:"path"`
	Methods    []string               `json:"methods"`
	Parameters []APIParameter         `json:"parameters"`
	Responses  map[string]interface{} `json:"responses"`
}

// APIParameter represents an API parameter
type APIParameter struct {
	Name     string                 `json:"name"`
	In       string                 `json:"in"`
	Schema   map[string]interface{} `json:"schema"`
	Required bool                   `json:"required"`
}

// GetProductList gets the list of all products
func (c *APIMetaClient) GetProductList() ([]ProductInfo, error) {
	url := fmt.Sprintf("%s/products.json", c.baseURL)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get product list: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API meta returned status %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	var products []ProductInfo
	if err := json.Unmarshal(body, &products); err != nil {
		return nil, fmt.Errorf("failed to unmarshal products: %w", err)
	}
	
	return products, nil
}

// GetServiceVersion gets the default version for a service
func (c *APIMetaClient) GetServiceVersion(service string) (string, error) {
	products, err := c.GetProductList()
	if err != nil {
		return "", err
	}
	
	serviceLower := strings.ToLower(service)
	for _, product := range products {
		if strings.ToLower(product.Code) == serviceLower {
			return product.DefaultVersion, nil
		}
	}
	
	return "", fmt.Errorf("service not found: %s", service)
}

// GetServiceStyle gets the style (RPC/ROA) for a service
func (c *APIMetaClient) GetServiceStyle(service string) (string, error) {
	products, err := c.GetProductList()
	if err != nil {
		return "RPC", nil // Default to RPC
	}
	
	serviceLower := strings.ToLower(service)
	for _, product := range products {
		if strings.ToLower(product.Code) == serviceLower {
			if product.Style != "" {
				return product.Style, nil
			}
			return "RPC", nil
		}
	}
	
	return "RPC", nil
}

// GetAPIInfo gets detailed information about a specific API
func (c *APIMetaClient) GetAPIInfo(service, api, version string) (*APIInfo, error) {
	url := fmt.Sprintf("%s/products/%s/versions/%s/apis/%s/api.json", 
		c.baseURL, service, version, api)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get API info: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API meta returned status %d for %s.%s", resp.StatusCode, service, api)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	var apiInfo APIInfo
	if err := json.Unmarshal(body, &apiInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal API info: %w", err)
	}
	
	return &apiInfo, nil
}

// GetAPIOverview gets overview of all APIs in a service
func (c *APIMetaClient) GetAPIOverview(service, version string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/products/%s/versions/%s/overview.json", 
		c.baseURL, service, version)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get API overview: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API meta returned status %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	var overview map[string]interface{}
	if err := json.Unmarshal(body, &overview); err != nil {
		return nil, fmt.Errorf("failed to unmarshal overview: %w", err)
	}
	
	return overview, nil
}

// GetAPIsInService gets the list of all APIs in a service
func (c *APIMetaClient) GetAPIsInService(service string) ([]string, error) {
	version, err := c.GetServiceVersion(service)
	if err != nil {
		return nil, err
	}
	
	overview, err := c.GetAPIOverview(service, version)
	if err != nil {
		return nil, err
	}
	
	apisInterface, ok := overview["apis"]
	if !ok {
		return []string{}, nil
	}
	
	apisMap, ok := apisInterface.(map[string]interface{})
	if !ok {
		return []string{}, nil
	}
	
	apis := make([]string, 0, len(apisMap))
	for apiName := range apisMap {
		apis = append(apis, apiName)
	}
	
	return apis, nil
}

// GetStandardServiceAndAPI gets the standard (correct case) service and API names
func (c *APIMetaClient) GetStandardServiceAndAPI(service, api, version string) (string, string, error) {
	products, err := c.GetProductList()
	if err != nil {
		return "", "", err
	}
	
	serviceLower := strings.ToLower(service)
	var standardService string
	for _, product := range products {
		if strings.ToLower(product.Code) == serviceLower {
			standardService = product.Code
			break
		}
	}
	
	if standardService == "" {
		return "", "", fmt.Errorf("service not found: %s", service)
	}
	
	if api == "" {
		return standardService, "", nil
	}
	
	// Get API overview to find standard API name
	overview, err := c.GetAPIOverview(standardService, version)
	if err != nil {
		return standardService, "", err
	}
	
	apisInterface, ok := overview["apis"]
	if !ok {
		return standardService, "", fmt.Errorf("no APIs found for service: %s", service)
	}
	
	apisMap, ok := apisInterface.(map[string]interface{})
	if !ok {
		return standardService, "", fmt.Errorf("invalid APIs format")
	}
	
	apiLower := strings.ToLower(api)
	for apiName := range apisMap {
		if strings.ToLower(apiName) == apiLower {
			return standardService, apiName, nil
		}
	}
	
	return standardService, "", fmt.Errorf("API not found: %s in service: %s", api, service)
}

