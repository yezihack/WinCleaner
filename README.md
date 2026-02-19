# 🧹 Win Cleaner

Windows 系统垃圾清理与优化工具，基于 Wails v2 构建，轻量高效。

## 功能

- **系统概览** — CPU、内存、磁盘使用率仪表盘，显卡信息检测（独显/核显自动识别）
- **垃圾清理** — 扫描 7 类系统垃圾（临时文件、Windows Update 缓存、缩略图、日志、浏览器缓存、回收站、预读取），支持展开查看文件列表，清理历史图表统计
- **内存优化** — 一键收缩进程工作集释放物理内存，优化历史趋势图、每日/月度释放量图表、优化前后对比
- **进程管理** — 进程列表按 CPU/内存排序，搜索过滤，结束进程
- **流量监控** — 实时网速、进程网络使用、每日/月度/年度流量趋势图、上传下载占比饼图
- **磁盘管理** — 所有分区空间概览、分区对比图表、大文件扫描（可选分区和最小文件大小）
- **实时状态栏** — 侧边栏底部实时显示 CPU、内存占比和网络速率

## 技术栈

| 层 | 技术 |
|---|---|
| 框架 | [Wails v2](https://wails.io) |
| 后端 | Go 1.24 |
| 前端 | Vue 3 + TypeScript + Element Plus + ECharts |
| 系统信息 | [gopsutil](https://github.com/shirou/gopsutil) |

程序完全本地运行，不上传任何数据。

## 环境要求

- Windows 10+
- Go 1.24+
- Node.js 18+
- [Wails CLI v2](https://wails.io/docs/gettingstarted/installation)

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

## 开发

```bash
cd win-cleaner

# 安装前端依赖
cd frontend && npm install && cd ..

# 开发模式（热重载）
wails dev
```

## 构建

```bash
wails build -platform windows/amd64
```

产物在 `build/bin/WinCleaner.exe`，单文件，解压即用。

## 项目结构

```
win-cleaner/
├── frontend/              # Vue3 前端
│   └── src/
│       ├── api/           # Wails 绑定调用封装
│       ├── views/         # 页面（Dashboard/Cleaner/Memory/Process/Network/Disk）
│       └── router/        # 路由
├── internal/
│   ├── app/               # Wails App 主结构与生命周期
│   ├── cleaner/           # 垃圾扫描、清理、历史记录
│   ├── memory/            # 内存优化、优化历史
│   ├── model/             # 数据模型
│   └── monitor/           # 系统监控（CPU/内存/磁盘/GPU/网络/进程）
├── pkg/winapi/            # Windows API 调用
├── build/                 # 构建资源（图标）
├── favicon_io/            # 应用图标源文件
├── main.go                # 入口
└── wails.json             # Wails 配置
```

## 数据存储

运行时数据保存在用户主目录下的 `.wincleaner/` 文件夹：

- `~/.wincleaner/clean_history.json` — 垃圾清理历史
- `~/.wincleaner/mem_opt_history.json` — 内存优化历史
- `~/.wincleaner/net_history.json` — 网络流量采样记录

Windows 下实际路径为 `C:\Users\<用户名>\.wincleaner\`。

## License

MIT
