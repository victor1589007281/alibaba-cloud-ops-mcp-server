package tools

import (
	"encoding/json"
	"fmt"
	
	"github.com/aliyun/alibaba-cloud-ops-mcp-server-go/pkg/mcp"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	oos "github.com/alibabacloud-go/oos-20190601/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

// RegisterOOSTools registers OOS-related tools
func RegisterOOSTools(r *Registry) error {
	// OOS_RunCommand
	r.RegisterTool(mcp.Tool{
		Name:        "OOS_RunCommand",
		Description: "批量在多台ECS实例上运行云助手命令，适用于需要同时管理多台ECS实例的场景，如应用程序管理和资源标记操作等。",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"Command": {
					Type:        "string",
					Description: "在ECS实例上执行的命令内容",
				},
				"InstanceIds": {
					Type:        "array",
					Description: "阿里云ECS实例ID列表",
					Items: &mcp.Property{
						Type: "string",
					},
				},
				"RegionId": {
					Type:        "string",
					Description: "阿里云地域ID",
					Default:     "cn-hangzhou",
				},
				"CommandType": {
					Type:        "string",
					Description: "命令类型: RunShellScript, RunPythonScript, RunPerlScript, RunBatScript, RunPowerShellScript",
					Default:     "RunShellScript",
				},
			},
			Required: []string{"Command", "InstanceIds"},
		},
	}, r.createOOSRunCommand())
	
	// OOS_StartInstances
	r.RegisterTool(mcp.Tool{
		Name:        "OOS_StartInstances",
		Description: "批量启动ECS实例，适用于需要同时管理和启动多台ECS实例的场景，例如应用部署和高可用性场景。",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"InstanceIds": {
					Type:        "array",
					Description: "阿里云ECS实例ID列表",
					Items: &mcp.Property{
						Type: "string",
					},
				},
				"RegionId": {
					Type:        "string",
					Description: "阿里云地域ID",
					Default:     "cn-hangzhou",
				},
			},
			Required: []string{"InstanceIds"},
		},
	}, r.createOOSStartInstances())
	
	// OOS_StopInstances
	r.RegisterTool(mcp.Tool{
		Name:        "OOS_StopInstances",
		Description: "批量停止ECS实例，适用于需要同时管理和停止多台ECS实例的场景。",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"InstanceIds": {
					Type:        "array",
					Description: "阿里云ECS实例ID列表",
					Items: &mcp.Property{
						Type: "string",
					},
				},
				"RegionId": {
					Type:        "string",
					Description: "阿里云地域ID",
					Default:     "cn-hangzhou",
				},
				"ForceStop": {
					Type:        "boolean",
					Description: "是否强制停止",
					Default:     false,
				},
			},
			Required: []string{"InstanceIds"},
		},
	}, r.createOOSStopInstances())
	
	// OOS_RebootInstances
	r.RegisterTool(mcp.Tool{
		Name:        "OOS_RebootInstances",
		Description: "批量重启ECS实例，适用于需要同时管理和重启多台ECS实例的场景。",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"InstanceIds": {
					Type:        "array",
					Description: "阿里云ECS实例ID列表",
					Items: &mcp.Property{
						Type: "string",
					},
				},
				"RegionId": {
					Type:        "string",
					Description: "阿里云地域ID",
					Default:     "cn-hangzhou",
				},
				"ForceStop": {
					Type:        "boolean",
					Description: "是否强制停止",
					Default:     false,
				},
			},
			Required: []string{"InstanceIds"},
		},
	}, r.createOOSRebootInstances())
	
	// Add more OOS tools...
	return nil
}

// createOOSClient creates an OOS client
func (r *Registry) createOOSClient(regionID string) (*oos.Client, error) {
	creds, err := r.GetCredentials()
	if err != nil {
		return nil, err
	}
	
	config := &openapi.Config{
		AccessKeyId:     tea.String(creds.AccessKeyID),
		AccessKeySecret: tea.String(creds.AccessKeySecret),
		Endpoint:        tea.String(fmt.Sprintf("oos.%s.aliyuncs.com", regionID)),
	}
	
	if creds.SecurityToken != "" {
		config.SecurityToken = tea.String(creds.SecurityToken)
	}
	
	return oos.NewClient(config)
}

// createOOSRunCommand creates the OOS_RunCommand tool function
func (r *Registry) createOOSRunCommand() ToolFunc {
	return func(args map[string]interface{}) (interface{}, error) {
		command, _ := args["Command"].(string)
		regionID, _ := args["RegionId"].(string)
		if regionID == "" {
			regionID = "cn-hangzhou"
		}
		commandType, _ := args["CommandType"].(string)
		if commandType == "" {
			commandType = "RunShellScript"
		}
		
		// Parse InstanceIds
		instanceIdsRaw, _ := args["InstanceIds"].([]interface{})
		instanceIds := make([]string, len(instanceIdsRaw))
		for i, v := range instanceIdsRaw {
			instanceIds[i] = fmt.Sprint(v)
		}
		
		client, err := r.createOOSClient(regionID)
		if err != nil {
			return nil, err
		}
		
		// Create execution parameters
		parameters := map[string]interface{}{
			"regionId":       regionID,
			"resourceType":   "ALIYUN::ECS::Instance",
			"targets": map[string]interface{}{
				"ResourceIds": instanceIds,
				"RegionId":    regionID,
				"Type":        "ResourceIds",
				"Parameters": map[string]string{
					"RegionId": regionID,
					"Status":   "Running",
				},
			},
			"commandType":    commandType,
			"commandContent": command,
		}
		
		return r.startExecutionSync(client, regionID, "ACS-ECS-BulkyRunCommand", parameters)
	}
}

// createOOSStartInstances creates the OOS_StartInstances tool function
func (r *Registry) createOOSStartInstances() ToolFunc {
	return func(args map[string]interface{}) (interface{}, error) {
		regionID, _ := args["RegionId"].(string)
		if regionID == "" {
			regionID = "cn-hangzhou"
		}
		
		instanceIdsRaw, _ := args["InstanceIds"].([]interface{})
		instanceIds := make([]string, len(instanceIdsRaw))
		for i, v := range instanceIdsRaw {
			instanceIds[i] = fmt.Sprint(v)
		}
		
		client, err := r.createOOSClient(regionID)
		if err != nil {
			return nil, err
		}
		
		parameters := map[string]interface{}{
			"regionId":     regionID,
			"resourceType": "ALIYUN::ECS::Instance",
			"targets": map[string]interface{}{
				"ResourceIds": instanceIds,
				"RegionId":    regionID,
				"Type":        "ResourceIds",
			},
		}
		
		return r.startExecutionSync(client, regionID, "ACS-ECS-BulkyStartInstances", parameters)
	}
}

// createOOSStopInstances creates the OOS_StopInstances tool function
func (r *Registry) createOOSStopInstances() ToolFunc {
	return func(args map[string]interface{}) (interface{}, error) {
		regionID, _ := args["RegionId"].(string)
		if regionID == "" {
			regionID = "cn-hangzhou"
		}
		forceStop, _ := args["ForceStop"].(bool)
		
		instanceIdsRaw, _ := args["InstanceIds"].([]interface{})
		instanceIds := make([]string, len(instanceIdsRaw))
		for i, v := range instanceIdsRaw {
			instanceIds[i] = fmt.Sprint(v)
		}
		
		client, err := r.createOOSClient(regionID)
		if err != nil {
			return nil, err
		}
		
		parameters := map[string]interface{}{
			"regionId":     regionID,
			"resourceType": "ALIYUN::ECS::Instance",
			"targets": map[string]interface{}{
				"ResourceIds": instanceIds,
				"RegionId":    regionID,
				"Type":        "ResourceIds",
			},
			"forceStop": forceStop,
		}
		
		return r.startExecutionSync(client, regionID, "ACS-ECS-BulkyStopInstances", parameters)
	}
}

// createOOSRebootInstances creates the OOS_RebootInstances tool function
func (r *Registry) createOOSRebootInstances() ToolFunc {
	return func(args map[string]interface{}) (interface{}, error) {
		regionID, _ := args["RegionId"].(string)
		if regionID == "" {
			regionID = "cn-hangzhou"
		}
		forceStop, _ := args["ForceStop"].(bool)
		
		instanceIdsRaw, _ := args["InstanceIds"].([]interface{})
		instanceIds := make([]string, len(instanceIdsRaw))
		for i, v := range instanceIdsRaw {
			instanceIds[i] = fmt.Sprint(v)
		}
		
		client, err := r.createOOSClient(regionID)
		if err != nil {
			return nil, err
		}
		
		parameters := map[string]interface{}{
			"regionId":     regionID,
			"resourceType": "ALIYUN::ECS::Instance",
			"targets": map[string]interface{}{
				"ResourceIds": instanceIds,
				"RegionId":    regionID,
				"Type":        "ResourceIds",
			},
			"forceStop": forceStop,
		}
		
		return r.startExecutionSync(client, regionID, "ACS-ECS-BulkyRebootInstances", parameters)
	}
}

// startExecutionSync starts an OOS execution and waits for completion
func (r *Registry) startExecutionSync(client *oos.Client, regionID, templateName string, parameters map[string]interface{}) (interface{}, error) {
	// Convert parameters to JSON string
	paramsJSON, err := json.Marshal(parameters)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal parameters: %w", err)
	}
	
	// Start execution
	startReq := &oos.StartExecutionRequest{
		RegionId:     tea.String(regionID),
		TemplateName: tea.String(templateName),
		Parameters:   tea.String(string(paramsJSON)),
	}
	
	startResp, err := client.StartExecution(startReq)
	if err != nil {
		return nil, fmt.Errorf("failed to start execution: %w", err)
	}
	
	executionID := tea.StringValue(startResp.Body.Execution.ExecutionId)
	r.logger.Printf("Started OOS execution: %s", executionID)
	
	// Poll for completion
	for {
		listReq := &oos.ListExecutionsRequest{
			RegionId:    tea.String(regionID),
			ExecutionId: tea.String(executionID),
		}
		
		listResp, err := client.ListExecutions(listReq)
		if err != nil {
			return nil, fmt.Errorf("failed to list executions: %w", err)
		}
		
		if len(listResp.Body.Executions) == 0 {
			return nil, fmt.Errorf("execution not found: %s", executionID)
		}
		
		execution := listResp.Body.Executions[0]
		status := tea.StringValue(execution.Status)
		
		r.logger.Printf("Execution %s status: %s", executionID, status)
		
		if status == "Failed" {
			statusMsg := tea.StringValue(execution.StatusMessage)
			return nil, fmt.Errorf("execution failed: %s", statusMsg)
		}
		
		if status == "Success" || status == "Cancelled" {
			return listResp.Body, nil
		}
		
		// Wait before next poll
		tea.Sleep(tea.Int(1000)) // 1 second
	}
}

