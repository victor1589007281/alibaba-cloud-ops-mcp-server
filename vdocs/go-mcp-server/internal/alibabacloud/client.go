package alibabacloud

import (
	"fmt"
	"strings"
	
	"github.com/aliyun/alibaba-cloud-ops-mcp-server-go/internal/config"
)

// ServiceClient represents a generic AlibabaCloud service client
type ServiceClient struct {
	Service     string
	RegionID    string
	Endpoint    string
	Credentials *Credentials
	UserAgent   string
}

// ClientConfig represents client configuration
type ClientConfig struct {
	Service     string
	RegionID    string
	Credentials *Credentials
	Settings    *config.Settings
}

// RegionEndpointServices are services that use region-specific endpoints
var RegionEndpointServices = map[string]bool{
	"ecs": true,
	"oos": true,
	"vpc": true,
	"slb": true,
}

// DoubleEndpointServices are services that use double endpoints
var DoubleEndpointServices = map[string][]string{
	"rds": {"cn-qingdao", "cn-beijing", "cn-hangzhou", "cn-shanghai", "cn-shenzhen", "cn-heyuan", "cn-guangzhou", "cn-hongkong"},
	"ess": {"cn-qingdao", "cn-beijing", "cn-hangzhou", "cn-shanghai", "cn-nanjing", "cn-shenzhen"},
	"dds": {"cn-qingdao", "cn-beijing", "cn-wulanchabu", "cn-hangzhou", "cn-shanghai", "cn-shenzhen", "cn-heyuan", "cn-guangzhou"},
	"r-kvstore": {"cn-qingdao", "cn-beijing", "cn-wulanchabu", "cn-hangzhou", "cn-shanghai", "cn-shenzhen", "cn-heyuan"},
}

// CentralServices are services that use central endpoints
var CentralServices = map[string]bool{
	"cbn": true,
	"ros": true,
	"ram": true,
}

// CentralServiceEndpoints defines special central service endpoints
var CentralServiceEndpoints = map[string]CentralEndpointConfig{
	"bssopenapi": {
		DomesticEndpoint:     "business.aliyuncs.com",
		InternationalEndpoint: "business.ap-southeast-1.aliyuncs.com",
		DomesticRegions: []string{"cn-qingdao", "cn-beijing", "cn-zhangjiakou", "cn-huhehaote", "cn-wulanchabu",
			"cn-hangzhou", "cn-shanghai", "cn-shenzhen", "cn-chengdu", "cn-hongkong"},
	},
}

// CentralEndpointConfig defines central endpoint configuration
type CentralEndpointConfig struct {
	DomesticEndpoint      string
	InternationalEndpoint string
	DomesticRegions       []string
}

// NewServiceClient creates a new service client
func NewServiceClient(cfg *ClientConfig) *ServiceClient {
	service := strings.ToLower(cfg.Service)
	regionID := strings.ToLower(cfg.RegionID)
	
	endpoint := getServiceEndpoint(service, regionID, cfg.Settings)
	
	return &ServiceClient{
		Service:     service,
		RegionID:    regionID,
		Endpoint:    endpoint,
		Credentials: cfg.Credentials,
		UserAgent:   "alibaba-cloud-ops-mcp-server-go/1.0.0",
	}
}

// getServiceEndpoint determines the endpoint for a service
func getServiceEndpoint(service, regionID string, settings *config.Settings) string {
	// Check central service endpoints first
	if centralConfig, ok := CentralServiceEndpoints[service]; ok {
		if settings.IsInternational() {
			return centralConfig.InternationalEndpoint
		}
		
		// Check if region is in domestic regions
		for _, region := range centralConfig.DomesticRegions {
			if region == regionID {
				return centralConfig.DomesticEndpoint
			}
		}
		
		if settings.IsDomestic() {
			return centralConfig.DomesticEndpoint
		}
		return centralConfig.InternationalEndpoint
	}
	
	// Check region endpoint services
	if RegionEndpointServices[service] {
		return fmt.Sprintf("%s.%s.aliyuncs.com", service, regionID)
	}
	
	// Check double endpoint services
	if regions, ok := DoubleEndpointServices[service]; ok {
		inCentralRegion := false
		for _, region := range regions {
			if region == regionID {
				inCentralRegion = true
				break
			}
		}
		
		if !inCentralRegion {
			return fmt.Sprintf("%s.%s.aliyuncs.com", service, regionID)
		}
		return fmt.Sprintf("%s.aliyuncs.com", service)
	}
	
	// Check central services
	if CentralServices[service] {
		return fmt.Sprintf("%s.aliyuncs.com", service)
	}
	
	// Default: region-specific endpoint
	return fmt.Sprintf("%s.%s.aliyuncs.com", service, regionID)
}

