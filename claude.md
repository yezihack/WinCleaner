# WinCleaner 项目规范

## 项目概述

WinCleaner 是一个基于 Wails 框架的 Windows 系统清理工具，采用 Go 后端 + Vue3 前端架构。

## 技术栈

- **后端**: Go 1.24+, Wails v2, GORM
- **前端**: Vue 3, TypeScript, Vite, Naive UI
- **平台**: Windows

## 项目结构

```
WinCleaner/
├── main.go              # Wails 应用入口
├── wails.json           # Wails 配置
├── internal/            # 内部包
│   ├── app/             # 应用核心
│   ├── cleaner/         # 清理模块（扫描、历史）
│   ├── memory/          # 内存优化模块
│   ├── monitor/         # 系统监控（进程、网络、端口、磁盘、GPU）
│   └── model/           # 数据模型
├── pkg/                 # 公共包
│   ├── datadir/         # 数据目录
│   └── winapi/          # Windows API 封装
└── frontend/            # Vue3 前端
    └── src/
        ├── api/         # 后端 API 调用
        ├── router/      # 路由配置
        └── views/       # 页面组件
```

## 开发规范

### Go 后端（遵循 Handler → Service → Repo → Model 分层）

- Handler: 处理请求参数校验、响应封装
- Service: 业务逻辑处理
- Repo: 数据访问层
- Model: 数据模型定义

### 前端 Vue3

- 使用 `<script setup lang="ts">` 语法
- Composition API
- TypeScript 严格模式

## API 约定

- RESTful 风格
- 响应格式: `{ "code": 0, "msg": "success", "data": {} }`
- 分页包含: records, total, page, page_size
