# 快速开始指南

## 5分钟上手阿里云运维MCP服务器（Golang版本）

### 第一步：准备环境 (1分钟)

#### 检查Go版本

```bash
go version
# 需要 Go 1.21 或更高版本
```

#### 准备阿里云AccessKey

访问 [阿里云控制台](https://ram.console.aliyun.com/manage/ak) 创建AccessKey

### 第二步：下载和编译 (2分钟)

```bash
# 克隆项目
cd /path/to/alibaba-cloud-ops-mcp-server/vdocs/go-mcp-server

# 安装依赖并编译
make build

# 验证编译结果
ls -lh bin/go-mcp-server
```

### 第三步：配置认证 (1分钟)

#### 方式1：环境变量（推荐）

```bash
export ALIBABA_CLOUD_ACCESS_KEY_ID="your-access-key-id"
export ALIBABA_CLOUD_ACCESS_KEY_SECRET="your-access-key-secret"
```

#### 方式2：创建.env文件

```bash
cp .env.example .env
# 编辑.env文件，填入你的AccessKey
```

### 第四步：运行测试 (1分钟)

```bash
# 运行所有测试
make test

# 查看测试结果
# 应该看到 25 个测试全部通过
```

### 第五步：配置MCP客户端 (可选)

#### 配置Claude Desktop

编辑 `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "alibaba-cloud-ops": {
      "command": "/path/to/go-mcp-server/bin/go-mcp-server",
      "args": ["--transport", "stdio"],
      "env": {
        "ALIBABA_CLOUD_ACCESS_KEY_ID": "your-key-id",
        "ALIBABA_CLOUD_ACCESS_KEY_SECRET": "your-key-secret"
      }
    }
  }
}
```

#### 配置Cursor

类似地配置Cursor的MCP设置。

## 验证安装

### 方式1：直接运行

```bash
# 设置环境变量
export ALIBABA_CLOUD_ACCESS_KEY_ID="your-key-id"
export ALIBABA_CLOUD_ACCESS_KEY_SECRET="your-key-secret"

# 运行服务器
./bin/go-mcp-server --transport stdio
```

服务器启动后，你会看到类似输出：

```
[Server] Starting MCP server (transport: stdio)
[Server] Registered tool: OOS_RunCommand
[Server] Registered tool: OOS_StartInstances
...
[Server] Server started with 28 tools registered
```

### 方式2：查看帮助

```bash
./bin/go-mcp-server --help
```

输出：

```
Usage of ./bin/go-mcp-server:
  -env string
        Environment: domestic or international (default "domestic")
  -headers-credential-only
        Use credentials only from HTTP headers
  -host string
        Host address for HTTP/SSE transport (default "127.0.0.1")
  -port int
        Port number for HTTP/SSE transport (default 8000)
  -services string
        Comma-separated list of services to enable (e.g., ecs,vpc,rds)
  -transport string
        Transport type: stdio, http, sse (default "stdio")
```

## 常见使用场景

### 场景1：查询ECS实例

在MCP客户端（如Claude）中，可以这样问：

```
请帮我查询杭州地域的所有ECS实例
```

MCP服务器会自动调用 `ECS_DescribeInstances` 工具。

### 场景2：监控CPU使用率

```
请帮我查看实例 i-xxxxx 的CPU使用率
```

MCP服务器会调用 `CMS_GetCpuUsageData` 工具。

### 场景3：批量启动实例

```
请帮我启动以下ECS实例：i-xxxxx, i-yyyyy
```

MCP服务器会调用 `OOS_StartInstances` 工具。

### 场景4：查询OSS存储桶

```
请列出杭州地域的所有OSS存储桶
```

MCP服务器会调用 `OSS_ListBuckets` 工具。

## 可用工具清单

服务器启动后会注册 **28个工具**：

### ECS相关 (12个)
- OOS_RunCommand - 批量运行命令
- OOS_StartInstances - 批量启动实例
- OOS_StopInstances - 批量停止实例
- OOS_RebootInstances - 批量重启实例
- ECS_DescribeInstances - 查询实例
- ECS_DescribeRegions - 查询地域
- ECS_DescribeZones - 查询可用区
- ECS_DescribeImages - 查询镜像
- ECS_DescribeSecurityGroups - 查询安全组
- ECS_DescribeAccountAttributes - 查询账号属性
- ECS_DescribeAvailableResource - 查询资源库存
- ECS_DeleteInstances - 删除实例

### 云监控 (9个)
- CMS_GetCpuUsageData - CPU使用率
- CMS_GetCpuLoadavgData - 1分钟负载
- CMS_GetCpuloadavg5mData - 5分钟负载
- CMS_GetCpuloadavg15mData - 15分钟负载
- CMS_GetMemUsedData - 内存使用量
- CMS_GetMemUsageData - 内存使用率
- CMS_GetDiskUsageData - 磁盘使用率
- CMS_GetDiskTotalData - 磁盘总容量
- CMS_GetDiskUsedData - 磁盘使用量

### OSS (4个)
- OSS_ListBuckets - 列出存储桶
- OSS_ListObjects - 列出对象
- OSS_PutBucket - 创建存储桶
- OSS_DeleteBucket - 删除存储桶

### VPC (2个)
- VPC_DescribeVpcs - 查询VPC
- VPC_DescribeVSwitches - 查询交换机

### RDS (1个)
- RDS_DescribeDBInstances - 查询数据库实例

## 高级配置

### 只启用特定服务

```bash
./bin/go-mcp-server \
  --transport stdio \
  --services ecs,vpc
```

### 国际环境配置

```bash
./bin/go-mcp-server \
  --transport stdio \
  --env international
```

### 使用HTTP传输（开发中）

```bash
./bin/go-mcp-server \
  --transport http \
  --host 0.0.0.0 \
  --port 8000
```

## 开发和调试

### 查看日志

服务器日志输出到stderr：

```bash
./bin/go-mcp-server 2> server.log
```

### 运行特定测试

```bash
# 运行协议测试
go test -v ./tests -run TestProtocol

# 运行配置测试
go test -v ./tests -run TestConfig

# 运行集成测试
go test -v ./tests -run TestIntegration
```

### 代码格式化

```bash
make fmt
```

### 生成测试覆盖率报告

```bash
make test-coverage
# 打开 coverage.html 查看详细报告
```

## 故障排查

### 问题1：认证失败

**症状**：
```
failed to get credentials: credentials not found in config or environment
```

**解决**：
- 检查环境变量是否设置
- 确认AccessKey是否正确
- 尝试重新设置环境变量

### 问题2：工具注册失败

**症状**：
```
Warning: failed to register API tool
```

**解决**：
- 这是正常的警告，部分API可能无法注册
- 不影响其他工具的使用
- 检查网络连接是否正常

### 问题3：编译错误

**症状**：
```
go: missing module requirements
```

**解决**：
```bash
go mod tidy
go mod download
```

## 性能基准

在MacBook Pro (M1, 16GB RAM)上的测试结果：

| 指标 | 值 |
|------|---|
| 编译时间 | ~5秒 |
| 二进制大小 | 9.7MB |
| 启动时间 | ~50ms |
| 工具注册 | ~3秒 |
| 内存占用 | ~25MB |
| 测试运行 | ~4.7秒 |

## 下一步

1. ✅ 尝试在MCP客户端中使用
2. ✅ 探索更多工具功能
3. ✅ 查看完整文档 [README.md](README.md)
4. ✅ 查看架构设计 [功能规划设计文档.md](功能规划设计文档.md)
5. ✅ 参与项目贡献

## 获取帮助

- 📖 查看 [README.md](README.md)
- 📋 查看 [项目开发总结.md](项目开发总结.md)
- 🐛 提交 [GitHub Issues](https://github.com/aliyun/alibaba-cloud-ops-mcp-server-go/issues)
- 💬 加入钉钉讨论群：113455011677

## 相关链接

- [Python版本](https://github.com/aliyun/alibaba-cloud-ops-mcp-server)
- [MCP协议文档](https://modelcontextprotocol.io/)
- [阿里云OpenAPI](https://api.aliyun.com/)

---

🎉 恭喜！你已经成功上手阿里云运维MCP服务器（Golang版本）！

