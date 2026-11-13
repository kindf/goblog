# GoBlog

一个使用 Go 编写的轻量级博客 Web 应用，提供用户注册与登录、文章浏览及文章管理等基础功能。项目采用 Go 标准库的 HTTP 服务，使用 Gorilla Mux 处理路由，并通过 GORM 访问 MySQL。

## 功能概览

- 用户注册、登录和退出，使用 Cookie 保存会话。
- 浏览文章列表和文章详情。
- 登录用户创建文章，并对自己的文章进行编辑或删除。
- 查看用户发布的文章。
- 启动时自动连接数据库并执行 GORM 数据模型迁移。

## 技术栈

- Go 1.16 或更高版本
- MySQL
- [Gorilla Mux](https://github.com/gorilla/mux)
- [GORM](https://gorm.io/)
- Go `html/template` 视图

## 快速开始

### 1. 准备数据库

创建一个 MySQL 数据库（默认名称为 `goblog`），并准备好数据库账号。应用首次启动时会自动迁移用户和文章数据表。

### 2. 配置环境变量

应用通过环境变量读取配置。数据库地址和账号按实际环境设置：

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| `APP_PORT` | `3000` | HTTP 服务端口 |
| `APP_ENV` | `production` | 运行环境 |
| `APP_DEBUG` | `false` | 调试模式 |
| `APP_KEY` | 项目内置默认值 | Cookie 会话加密密钥；部署时请设置为独立的随机密钥 |
| `DB_HOST` | `127.0.0.1` | MySQL 主机 |
| `DB_PORT` | `3306` | MySQL 端口 |
| `DB_DATABASE` | `goblog` | 数据库名称 |
| `DB_USERNAME` | 空 | 数据库用户名 |
| `DB_PASSWORD` | 空 | 数据库密码 |
| `DB_MAX_IDLE_CONNECTIONS` | `100` | 最大空闲连接数 |
| `DB_MAX_OPEN_CONNECTIONS` | `25` | 最大打开连接数 |
| `DB_MAX_LIFE_SECONDS` | `300` | 连接最大生命周期（秒） |
| `SESSION_NAME` | `goblog-session` | 会话 Cookie 名称 |

例如，在 Linux/macOS shell 中：

```bash
export DB_DATABASE=goblog
export DB_USERNAME=your_user
export DB_PASSWORD=your_password
export APP_KEY='请替换为安全的随机密钥'
go run .
```

启动后访问 <http://localhost:3000>。如果修改了 `APP_PORT`，请使用相应端口。

### 构建

```bash
go build -o goblog .
./goblog
```

## 项目结构

```text
.
├── app/
│   ├── http/          # HTTP 控制器与中间件
│   ├── models/        # 用户、文章等业务模型
│   ├── policies/     # 资源访问策略
│   └── requests/     # 表单校验
├── bootstrap/         # 数据库迁移与路由初始化
├── config/            # 应用、数据库和会话配置
├── pkg/               # 配置、路由、视图等通用组件
├── public/            # 公共资源
├── static/            # 静态文件
└── main.go            # 应用入口
```

## 开发说明

依赖由 Go Modules 管理。下载依赖并运行应用：

```bash
go mod download
go run .
```

请确保 MySQL 已启动，且配置的数据库及账号具备相应访问权限。
