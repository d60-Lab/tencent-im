# 运行测试指南

本项目的测试需要有效的腾讯云 IM 凭证。为了安全起见，凭证通过环境变量配置。

## 快速开始

### 1. 配置凭证

#### 方式一：使用 .env 文件（推荐）

```bash
# 复制示例文件
cp .env.example .env

# 编辑 .env 文件，填入您的实际凭证
vim .env
```

#### 方式二：直接设置环境变量

```bash
export TIM_APP_ID=1400000000
export TIM_APP_SECRET=your_app_secret_here
export TIM_USER_ID=administrator
export TIM_BASE_URL=https://console.tim.qq.com
```

#### 方式三：在命令行临时设置

```bash
TIM_APP_ID=1400000000 TIM_APP_SECRET=your_secret go test ./...
```

### 2. 运行测试

```bash
# 运行所有测试
go test ./...

# 运行单个包的测试
go test .

# 运行特定测试
go test -run TestIm_Account_ImportAccount .

# 跳过耗时的集成测试（使用 -short 标志）
go test -short ./...

# 查看详细输出
go test -v ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 环境变量说明

| 变量名 | 必填 | 默认值 | 说明 |
|--------|------|--------|------|
| `TIM_APP_ID` | ✅ | 无 | 应用 SDK AppID |
| `TIM_APP_SECRET` | ✅ | 无 | 应用密钥 |
| `TIM_USER_ID` | ❌ | `administrator` | 管理员账号 |
| `TIM_BASE_URL` | ❌ | `https://console.tim.qq.com` | API 基础 URL |

## 获取凭证

1. 登录 [腾讯云 IM 控制台](https://console.cloud.tencent.com/im)
2. 选择您的应用或创建新应用
3. 在应用详情页面获取 **SDKAppID** 和 **密钥**

## 不同地区的 API 域名

根据您的应用创建地区，选择对应的 `TIM_BASE_URL`：

- **中国**：`https://console.tim.qq.com`
- **新加坡**：`https://adminapisgp.im.qcloud.com`
- **韩国**：`https://adminapikr.im.qcloud.com`
- **德国**：`https://adminapiger.im.qcloud.com`
- **印度**：`https://adminapiind.im.qcloud.com`

## 安全注意事项

⚠️ **重要**：

1. ❌ **不要**将 `.env` 文件提交到代码仓库
2. ❌ **不要**在代码中硬编码凭证
3. ✅ **使用**环境变量或密钥管理服务
4. ✅ **确保** `.env` 已在 `.gitignore` 中

## CI/CD 环境配置

在 GitHub Actions 或其他 CI/CD 平台中，通过 Secrets 配置环境变量：

**GitHub Actions 示例**：

```yaml
- name: Run tests
  env:
    TIM_APP_ID: ${{ secrets.TIM_APP_ID }}
    TIM_APP_SECRET: ${{ secrets.TIM_APP_SECRET }}
  run: go test ./...
```

## 故障排查

### 测试失败：code: 60026

**错误信息**：`sdkappid illegal or not match domain`

**解决方案**：
- 检查 `TIM_APP_ID` 是否正确
- 检查 `TIM_BASE_URL` 是否与应用地区匹配

### 测试失败：code: 90012

**错误信息**：`invalid To_Account`

**解决方案**：
- 确保测试账号已导入 IM 系统
- 运行账号导入测试：`go test -run TestIm_Account_ImportAccount`

### 环境变量未读取

**解决方案**：
- 确保已设置环境变量
- 使用 `echo $TIM_APP_ID` 检查变量是否存在
- 在同一终端会话中运行测试

## 更多信息

- [腾讯云 IM 文档](https://cloud.tencent.com/document/product/269)
- [项目 README](README.md)
- [贡献指南](CONTRIBUTING.md)
