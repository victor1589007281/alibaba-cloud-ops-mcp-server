package tools

import (
	"encoding/json"
	"fmt"
	
	"github.com/aliyun/alibaba-cloud-ops-mcp-server-go/pkg/mcp"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	cms "github.com/alibabacloud-go/cms-20190101/v9/client"
	"github.com/alibabacloud-go/tea/tea"
)

// RegisterCMSTools registers CMS-related tools
func RegisterCMSTools(r *Registry) error {
	metrics := []struct {
		name        string
		description string
		metricName  string
	}{
		{"CMS_GetCpuUsageData", "获取ECS实例的CPU使用率数据", "cpu_total"},
		{"CMS_GetCpuLoadavgData", "获取CPU一分钟平均负载指标数据", "load_1m"},
		{"CMS_GetCpuloadavg5mData", "获取CPU五分钟平均负载指标数据", "load_5m"},
		{"CMS_GetCpuloadavg15mData", "获取CPU十五分钟平均负载指标数据", "load_15m"},
		{"CMS_GetMemUsedData", "获取内存使用量指标数据", "memory_usedspace"},
		{"CMS_GetMemUsageData", "获取内存利用率指标数据", "memory_usedutilization"},
		{"CMS_GetDiskUsageData", "获取磁盘利用率指标数据", "diskusage_utilization"},
		{"CMS_GetDiskTotalData", "获取磁盘分区总容量指标数据", "diskusage_total"},
		{"CMS_GetDiskUsedData", "获取磁盘分区使用量指标数据", "diskusage_used"},
	}
	
	for _, metric := range metrics {
		metricName := metric.metricName // Capture for closure
		r.RegisterTool(mcp.Tool{
			Name:        metric.name,
			Description: metric.description,
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
		}, r.createCMSGetMetricData(metricName))
	}
	
	return nil
}

// createCMSClient creates a CMS client
func (r *Registry) createCMSClient(regionID string) (*cms.Client, error) {
	creds, err := r.GetCredentials()
	if err != nil {
		return nil, err
	}
	
	config := &openapi.Config{
		AccessKeyId:     tea.String(creds.AccessKeyID),
		AccessKeySecret: tea.String(creds.AccessKeySecret),
		Endpoint:        tea.String(fmt.Sprintf("metrics.%s.aliyuncs.com", regionID)),
	}
	
	if creds.SecurityToken != "" {
		config.SecurityToken = tea.String(creds.SecurityToken)
	}
	
	return cms.NewClient(config)
}

// createCMSGetMetricData creates a CMS metric data retrieval tool function
func (r *Registry) createCMSGetMetricData(metricName string) ToolFunc {
	return func(args map[string]interface{}) (interface{}, error) {
		regionID, _ := args["RegionId"].(string)
		if regionID == "" {
			regionID = "cn-hangzhou"
		}
		
		instanceIdsRaw, _ := args["InstanceIds"].([]interface{})
		if len(instanceIdsRaw) == 0 {
			return nil, fmt.Errorf("InstanceIds is required")
		}
		
		// Build dimensions
		dimensions := make([]map[string]string, 0, len(instanceIdsRaw))
		for _, v := range instanceIdsRaw {
			dimensions = append(dimensions, map[string]string{
				"instanceId": fmt.Sprint(v),
			})
		}
		
		dimensionsJSON, err := json.Marshal(dimensions)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal dimensions: %w", err)
		}
		
		client, err := r.createCMSClient(regionID)
		if err != nil {
			return nil, err
		}
		
		req := &cms.DescribeMetricLastRequest{
			Namespace:  tea.String("acs_ecs_dashboard"),
			MetricName: tea.String(metricName),
			Dimensions: tea.String(string(dimensionsJSON)),
		}
		
		resp, err := client.DescribeMetricLast(req)
		if err != nil {
			return nil, fmt.Errorf("failed to get metric data: %w", err)
		}
		
		r.logger.Printf("CMS metric response: %s", tea.StringValue(resp.Body.Datapoints))
		
		return map[string]interface{}{
			"code":       tea.StringValue(resp.Body.Code),
			"message":    tea.StringValue(resp.Body.Message),
			"datapoints": tea.StringValue(resp.Body.Datapoints),
			"period":     tea.StringValue(resp.Body.Period),
		}, nil
	}
}

