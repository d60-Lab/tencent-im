# 贡献指南

感谢您对 Tencent IM SDK for Go 的关注！我们欢迎任何形式的贡献。

## 🤝 如何贡献

### 报告 Bug

如果您发现了 bug，请通过 [GitHub Issues](https://github.com/d60-Lab/tencent-im/issues) 提交问题，并包含以下信息：

- **问题描述**：清晰简洁的描述
- **复现步骤**：详细的复现步骤
- **期望行为**：您期望发生什么
- **实际行为**：实际发生了什么
- **环境信息**：Go 版本、操作系统等
- **代码示例**：最小化的可复现代码

### 提交功能请求

如果您有新功能的想法，请先创建一个 [Discussion](https://github.com/d60-Lab/tencent-im/discussions) 或 Issue 讨论：

- **功能描述**：描述您想要的功能
- **使用场景**：为什么需要这个功能
- **建议方案**：如果有的话，描述您的实现思路

### 提交代码

1. **Fork 仓库**
   ```bash
   git clone https://github.com/your-username/tencent-im.git
   cd tencent-im
   ```

2. **创建分支**
   ```bash
   git checkout -b feature/your-feature-name
   # 或
   git checkout -b fix/your-bug-fix
   ```

3. **安装开发工具**
   ```bash
   # 安装 pre-commit hooks
   ./scripts/install-hooks.sh
   
   # 安装代码检查工具（可选但推荐）
   go install golang.org/x/tools/cmd/goimports@latest
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

4. **编写代码**
   - 遵循 Go 代码规范
   - 为新功能添加测试
   - 为公共API添加文档注释
   - 确保代码通过所有测试

5. **提交代码**
   ```bash
   git add .
   git commit -m "feat: add new feature"  # 使用 Conventional Commits 规范
   ```

   **Commit 消息规范**：
   - `feat:` 新功能
   - `fix:` Bug 修复
   - `docs:` 文档更新
   - `style:` 代码格式调整
   - `refactor:` 代码重构
   - `test:` 测试相关
   - `chore:` 构建/工具链相关

6. **推送并创建 PR**
   ```bash
   git push origin feature/your-feature-name
   ```
   
   然后在 GitHub 上创建 Pull Request。

## 📝 开发规范

### 代码风格

- 使用 `gofmt` 格式化代码
- 使用 `goimports` 管理导入
- 遵循 [Effective Go](https://go.dev/doc/effective_go) 指南
- 变量命名使用驼峰命名法
- 导出的标识符必须有文档注释

### 测试规范

- 为新功能编写单元测试
- 测试覆盖率应保持在 70% 以上
- 测试文件命名为 `*_test.go`
- 使用表驱动测试（table-driven tests）

示例：
```go
func TestNewClient(t *testing.T) {
    tests := []struct {
        name    string
        opt     *Options
        wantErr bool
    }{
        {
            name: "valid options",
            opt: &Options{
                AppId:     1400000000,
                AppSecret: "secret",
                UserId:    "admin",
            },
            wantErr: false,
        },
        // 更多测试用例...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := NewClient(tt.opt)
            if (got == nil) != tt.wantErr {
                t.Errorf("NewClient() error = %v, wantErr %v", got, tt.wantErr)
            }
        })
    }
}
```

### 文档规范

- 所有导出的函数、类型、常量都必须有文档注释
- 文档注释以标识符名称开头
- 复杂功能需要提供示例代码

示例：
```go
// NewIM creates a new Tencent IM client instance.
// It requires AppId, AppSecret and admin UserId.
//
// Example:
//     tim := im.NewIM(&im.Options{
//         AppId:     1400000000,
//         AppSecret: "your-secret",
//         UserId:    "administrator",
//     })
func NewIM(opt *Options) IM {
    // ...
}
```

## 🔍 代码审查流程

提交的 PR 将经过以下检查：

1. **自动化检查**
   - CI 测试必须通过
   - 代码格式检查必须通过
   - 安全扫描必须通过
   - 代码覆盖率不能降低

2. **人工审查**
   - 代码质量和可维护性
   - 是否符合项目架构
   - 文档和测试是否完善
   - 是否有破坏性变更

3. **合并要求**
   - 至少一位维护者 approve
   - 所有讨论已解决
   - CI 全部通过

## 📞 获取帮助

- **文档**：查看 [README.md](README.md) 和代码注释
- **Discussion**：在 [GitHub Discussions](https://github.com/d60-Lab/tencent-im/discussions) 提问
- **Issue**：查看现有 [Issues](https://github.com/d60-Lab/tencent-im/issues) 或创建新的

## 📜 行为准则

参与本项目即表示您同意遵守我们的 [行为准则](CODE_OF_CONDUCT.md)。

## 🙏 致谢

感谢所有为本项目做出贡献的开发者！
