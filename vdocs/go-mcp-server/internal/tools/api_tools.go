package tools

import (
	"encoding/json"
	"fmt"
	"strings"
	
	"github.com/aliyun/alibaba-cloud-ops-mcp-server-go/internal/alibabacloud"
	"github.com/aliyun/alibaba-cloud-ops-mcp-server-go/internal/config"
	"github.com/aliyun/alibaba-cloud-ops-mcp-server-go/pkg/mcp"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

// APIToolsConfig defines which APIs to expose as tools
var APIToolsConfig = map[string][]string{
	"ecs": {
		"DescribeInstances",
		"DescribeRegions",
		"DescribeZones",
		"DescribeAccountAttributes",
		"DescribeAvailableResource",
		"DescribeImages",
		"DescribeSecurityGroups",
		"DeleteInstances",
	},
	"vpc": {
		"DescribeVpcs",
		"DescribeVSwitches",
	},
	"rds": {
		"DescribeDBInstances",
	},
}

// RegisterAPITools registers dynamic API tools
func RegisterAPITools(r *Registry) error {
	for service, apis := range APIToolsConfig {
		for _, api := range apis {
			if err := r.registerAPITool(service, api); err != nil {
				r.logger.Printf("Warning: failed to register API tool %s.%s: %v", service, api, err)
				// Continue registering other tools even if one fails
			}
		}
	}
	return nil
}

// registerAPITool registers a single API tool
func (r *Registry) registerAPITool(service, api string) error {
	apiMeta := r.GetAPIMetaClient()
	
	// Get API version
	version, err := apiMeta.GetServiceVersion(service)
	if err != nil {
		return fmt.Errorf("failed to get service version: %w", err)
	}
	
	// Get standard service and API names
	standardService, standardAPI, err := apiMeta.GetStandardServiceAndAPI(service, api, version)
	if err != nil {
		return fmt.Errorf("failed to get standard names: %w", err)
	}
	
	// Get API info
	apiInfo, err := apiMeta.GetAPIInfo(standardService, standardAPI, version)
	if err != nil {
		return fmt.Errorf("failed to get API info: %w", err)
	}
	
	// Build tool schema
	tool := mcp.Tool{
		Name:        fmt.Sprintf("%s_%s", strings.ToUpper(service), standardAPI),
		Description: apiInfo.Summary,
		InputSchema: r.buildInputSchema(apiInfo),
	}
	
	// Create tool function
	toolFunc := r.createAPIToolFunc(standardService, standardAPI, version)
	
	// Register the tool
	r.RegisterTool(tool, toolFunc)
	
	return nil
}

// buildInputSchema builds MCP input schema from API parameters
func (r *Registry) buildInputSchema(apiInfo *alibabacloud.APIInfo) mcp.InputSchema {
	schema := mcp.InputSchema{
		Type:       "object",
		Properties: make(map[string]mcp.Property),
		Required:   []string{},
	}
	
	for _, param := range apiInfo.Parameters {
		// Skip parameters with dots (complex nested parameters)
		if strings.Contains(param.Name, ".") {
			continue
		}
		
		paramType := "string"
		if schemaType, ok := param.Schema["type"].(string); ok {
			paramType = schemaType
		}
		
		description := ""
		if desc, ok := param.Schema["description"].(string); ok {
			description = desc
		}
		
		example := ""
		if ex, ok := param.Schema["example"].(string); ok {
			example = ex
		}
		
		if example != "" {
			description = fmt.Sprintf("%s 参数示例: %s", description, example)
		}
		
		property := mcp.Property{
			Type:        paramType,
			Description: description,
		}
		
		schema.Properties[param.Name] = property
		
		if param.Required {
			schema.Required = append(schema.Required, param.Name)
		}
	}
	
	// Add RegionId if not present
	if _, exists := schema.Properties["RegionId"]; !exists {
		schema.Properties["RegionId"] = mcp.Property{
			Type:        "string",
			Description: "地域ID",
			Default:     "cn-hangzhou",
		}
	}
	
	return schema
}

// createAPIToolFunc creates a tool function for calling an API
func (r *Registry) createAPIToolFunc(service, api, version string) ToolFunc {
	return func(args map[string]interface{}) (interface{}, error) {
		// Get region ID
		regionID, _ := args["RegionId"].(string)
		if regionID == "" {
			regionID = "cn-hangzhou"
		}
		
		// Get credentials
		creds, err := r.GetCredentials()
		if err != nil {
			return nil, err
		}
		
		// Get service endpoint
		clientCfg := &alibabacloud.ClientConfig{
			Service:     service,
			RegionID:    regionID,
			Credentials: creds,
			Settings:    config.GetSettings(),
		}
		serviceClient := alibabacloud.NewServiceClient(clientCfg)
		
		// Get API style
		style, err := r.apiMeta.GetServiceStyle(service)
		if err != nil {
			style = "RPC"
		}
		
		// Get API info for method
		apiInfo, err := r.apiMeta.GetAPIInfo(service, api, version)
		if err != nil {
			return nil, fmt.Errorf("failed to get API info: %w", err)
		}
		
		method := "POST"
		if len(apiInfo.Methods) > 0 {
			if strings.ToLower(apiInfo.Methods[0]) == "post" {
				method = "POST"
			} else {
				method = "GET"
			}
		}
		
		// Clean up parameters (remove nil values)
		params := make(map[string]interface{})
		for k, v := range args {
			if v != nil {
				params[k] = v
			}
		}
		
		// Call the API
		return r.callOpenAPI(serviceClient, api, version, method, style, apiInfo.Path, params)
	}
}

// callOpenAPI calls an AlibabaCloud OpenAPI
func (r *Registry) callOpenAPI(
	serviceClient *alibabacloud.ServiceClient,
	action, version, method, style, path string,
	params map[string]interface{},
) (interface{}, error) {
	// Create OpenAPI client
	config := &openapi.Config{
		AccessKeyId:     tea.String(serviceClient.Credentials.AccessKeyID),
		AccessKeySecret: tea.String(serviceClient.Credentials.AccessKeySecret),
		Endpoint:        tea.String(serviceClient.Endpoint),
		UserAgent:       tea.String(serviceClient.UserAgent),
	}
	
	if serviceClient.Credentials.SecurityToken != "" {
		config.SecurityToken = tea.String(serviceClient.Credentials.SecurityToken)
	}
	
	client, err := openapi.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAPI client: %w", err)
	}
	
	// Convert params to query format
	query := make(map[string]*string)
	for k, v := range params {
		query[k] = tea.String(fmt.Sprint(v))
	}
	
	// Create OpenAPI request
	request := &openapi.OpenApiRequest{
		Query: query,
	}
	
	// Create params
	apiParams := &openapi.Params{
		Action:      tea.String(action),
		Version:     tea.String(version),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String(path),
		Method:      tea.String(method),
		AuthType:    tea.String("AK"),
		Style:       tea.String(style),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	
	// Call API
	runtime := &util.RuntimeOptions{}
	response, err := client.CallApi(apiParams, request, runtime)
	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}
	
	// Convert response to map
	body := response["body"]
	
	// If body is a map, return it directly
	if bodyMap, ok := body.(map[string]interface{}); ok {
		return bodyMap, nil
	}
	
	// Try to parse as JSON
	if bodyBytes, err := json.Marshal(body); err == nil {
		var result map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &result); err == nil {
			return result, nil
		}
	}
	
	return body, nil
}

