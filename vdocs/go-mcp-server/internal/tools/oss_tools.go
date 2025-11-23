package tools

import (
	"fmt"
	
	"github.com/aliyun/alibaba-cloud-ops-mcp-server-go/pkg/mcp"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// RegisterOSSTools registers OSS-related tools
func RegisterOSSTools(r *Registry) error {
	// OSS_ListBuckets
	r.RegisterTool(mcp.Tool{
		Name:        "OSS_ListBuckets",
		Description: "列出指定区域的所有OSS存储空间。",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"RegionId": {
					Type:        "string",
					Description: "阿里云地域ID",
					Default:     "cn-hangzhou",
				},
				"Prefix": {
					Type:        "string",
					Description: "OSS存储空间名称前缀",
					Default:     "",
				},
			},
		},
	}, r.createOSSListBuckets())
	
	// OSS_ListObjects
	r.RegisterTool(mcp.Tool{
		Name:        "OSS_ListObjects",
		Description: "获取指定OSS存储空间中的所有文件信息。",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"BucketName": {
					Type:        "string",
					Description: "OSS存储空间名称",
				},
				"RegionId": {
					Type:        "string",
					Description: "阿里云地域ID",
					Default:     "cn-hangzhou",
				},
				"Prefix": {
					Type:        "string",
					Description: "对象名称前缀",
					Default:     "",
				},
			},
			Required: []string{"BucketName"},
		},
	}, r.createOSSListObjects())
	
	// OSS_PutBucket
	r.RegisterTool(mcp.Tool{
		Name:        "OSS_PutBucket",
		Description: "创建一个新的OSS存储空间。",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"BucketName": {
					Type:        "string",
					Description: "OSS存储空间名称",
				},
				"RegionId": {
					Type:        "string",
					Description: "阿里云地域ID",
					Default:     "cn-hangzhou",
				},
				"StorageClass": {
					Type:        "string",
					Description: "存储类型: Standard, IA, Archive, ColdArchive, DeepColdArchive",
					Default:     "Standard",
				},
				"DataRedundancyType": {
					Type:        "string",
					Description: "数据容灾类型: LRS, ZRS",
					Default:     "LRS",
				},
			},
			Required: []string{"BucketName"},
		},
	}, r.createOSSPutBucket())
	
	// OSS_DeleteBucket
	r.RegisterTool(mcp.Tool{
		Name:        "OSS_DeleteBucket",
		Description: "删除指定的OSS存储空间。",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"BucketName": {
					Type:        "string",
					Description: "OSS存储空间名称",
				},
				"RegionId": {
					Type:        "string",
					Description: "阿里云地域ID",
					Default:     "cn-hangzhou",
				},
			},
			Required: []string{"BucketName"},
		},
	}, r.createOSSDeleteBucket())
	
	return nil
}

// createOSSClient creates an OSS client
func (r *Registry) createOSSClient(regionID string) (*oss.Client, error) {
	creds, err := r.GetCredentials()
	if err != nil {
		return nil, err
	}
	
	endpoint := fmt.Sprintf("oss-%s.aliyuncs.com", regionID)
	
	if creds.SecurityToken != "" {
		return oss.New(endpoint, creds.AccessKeyID, creds.AccessKeySecret, oss.SecurityToken(creds.SecurityToken))
	}
	
	return oss.New(endpoint, creds.AccessKeyID, creds.AccessKeySecret)
}

// createOSSListBuckets creates the OSS_ListBuckets tool function
func (r *Registry) createOSSListBuckets() ToolFunc {
	return func(args map[string]interface{}) (interface{}, error) {
		regionID, _ := args["RegionId"].(string)
		if regionID == "" {
			regionID = "cn-hangzhou"
		}
		prefix, _ := args["Prefix"].(string)
		
		client, err := r.createOSSClient(regionID)
		if err != nil {
			return nil, err
		}
		
		marker := ""
		var buckets []map[string]interface{}
		
		for {
			lsRes, err := client.ListBuckets(
				oss.Prefix(prefix),
				oss.Marker(marker),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to list buckets: %w", err)
			}
			
			for _, bucket := range lsRes.Buckets {
				buckets = append(buckets, map[string]interface{}{
					"name":         bucket.Name,
					"location":     bucket.Location,
					"creationDate": bucket.CreationDate,
					"storageClass": bucket.StorageClass,
				})
			}
			
			if !lsRes.IsTruncated {
				break
			}
			marker = lsRes.NextMarker
		}
		
		return buckets, nil
	}
}

// createOSSListObjects creates the OSS_ListObjects tool function
func (r *Registry) createOSSListObjects() ToolFunc {
	return func(args map[string]interface{}) (interface{}, error) {
		bucketName, _ := args["BucketName"].(string)
		if bucketName == "" {
			return nil, fmt.Errorf("BucketName is required")
		}
		
		regionID, _ := args["RegionId"].(string)
		if regionID == "" {
			regionID = "cn-hangzhou"
		}
		prefix, _ := args["Prefix"].(string)
		
		client, err := r.createOSSClient(regionID)
		if err != nil {
			return nil, err
		}
		
		bucket, err := client.Bucket(bucketName)
		if err != nil {
			return nil, fmt.Errorf("failed to get bucket: %w", err)
		}
		
		marker := ""
		var objects []map[string]interface{}
		
		for {
			lsRes, err := bucket.ListObjects(
				oss.Prefix(prefix),
				oss.Marker(marker),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to list objects: %w", err)
			}
			
			for _, object := range lsRes.Objects {
				objects = append(objects, map[string]interface{}{
					"key":          object.Key,
					"size":         object.Size,
					"lastModified": object.LastModified,
					"etag":         object.ETag,
					"storageClass": object.StorageClass,
				})
			}
			
			if !lsRes.IsTruncated {
				break
			}
			marker = lsRes.NextMarker
		}
		
		return objects, nil
	}
}

// createOSSPutBucket creates the OSS_PutBucket tool function
func (r *Registry) createOSSPutBucket() ToolFunc {
	return func(args map[string]interface{}) (interface{}, error) {
		bucketName, _ := args["BucketName"].(string)
		if bucketName == "" {
			return nil, fmt.Errorf("BucketName is required")
		}
		
		regionID, _ := args["RegionId"].(string)
		if regionID == "" {
			regionID = "cn-hangzhou"
		}
		storageClass, _ := args["StorageClass"].(string)
		if storageClass == "" {
			storageClass = "Standard"
		}
		dataRedundancy, _ := args["DataRedundancyType"].(string)
		if dataRedundancy == "" {
			dataRedundancy = "LRS"
		}
		
		client, err := r.createOSSClient(regionID)
		if err != nil {
			return nil, err
		}
		
		var storageClassType oss.StorageClassType
		switch storageClass {
		case "IA":
			storageClassType = oss.StorageIA
		case "Archive":
			storageClassType = oss.StorageArchive
		case "ColdArchive":
			storageClassType = oss.StorageColdArchive
		default:
			storageClassType = oss.StorageStandard
		}
		
		var redundancyType oss.DataRedundancyType
		if dataRedundancy == "ZRS" {
			redundancyType = oss.RedundancyZRS
		} else {
			redundancyType = oss.RedundancyLRS
		}
		
		err = client.CreateBucket(
			bucketName,
			oss.StorageClass(storageClassType),
			oss.RedundancyType(redundancyType),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
		
		return map[string]interface{}{
			"status":  "success",
			"message": fmt.Sprintf("Bucket %s created successfully", bucketName),
		}, nil
	}
}

// createOSSDeleteBucket creates the OSS_DeleteBucket tool function
func (r *Registry) createOSSDeleteBucket() ToolFunc {
	return func(args map[string]interface{}) (interface{}, error) {
		bucketName, _ := args["BucketName"].(string)
		if bucketName == "" {
			return nil, fmt.Errorf("BucketName is required")
		}
		
		regionID, _ := args["RegionId"].(string)
		if regionID == "" {
			regionID = "cn-hangzhou"
		}
		
		client, err := r.createOSSClient(regionID)
		if err != nil {
			return nil, err
		}
		
		err = client.DeleteBucket(bucketName)
		if err != nil {
			return nil, fmt.Errorf("failed to delete bucket: %w", err)
		}
		
		return map[string]interface{}{
			"status":  "success",
			"message": fmt.Sprintf("Bucket %s deleted successfully", bucketName),
		}, nil
	}
}

