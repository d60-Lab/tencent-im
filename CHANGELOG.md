# 更新日志

本项目的所有重要更改都将记录在此文件中。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，本项目遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [Unreleased]

### Added
- 配置 git-chglog 自动生成 CHANGELOG
- 添加 GitHub Actions 自动发布流程
- 添加 LICENSE 文件（MIT License）
- 添加 CONTRIBUTING.md 贡献指南
- 添加 CODE_OF_CONDUCT.md 行为准则
- 添加 pre-commit hooks 配置
- 添加日志接口支持
- 添加完善的配置选项（超时、重试、日志等）

### Changed
- 升级 Go 版本从 1.16 到 1.21
- 更新 API 域名为官方最新推荐域名
- 支持用户自定义 API 域名配置
- 替换 github.com/dobyte/http 为标准库 net/http
- 更新模块路径为 github.com/d60-Lab/tencent-im

### Fixed
- 修复 HTTP 客户端自动切换备用域名的逻辑
- 补充缺失的群组管理 API（计数器、自定义属性、禁言等）
- 补充缺失的会话管理 API
- 补充缺失的消息修改 API

## [0.1.0] - 2021

### Added
- 初始版本发布
- 账号管理 API
- 资料管理 API
- 关系链管理 API
- 群组管理 API
- 单聊消息 API
- 群聊消息 API
- 全员推送 API
- 运营管理 API
- 全局禁言 API
- 回调事件处理

[Unreleased]: https://github.com/d60-Lab/tencent-im/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/d60-Lab/tencent-im/releases/tag/v0.1.0
