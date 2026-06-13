# 校车预定平台 (Bus Reservation Platform)

本项目是一个基于 Go 语言开发的高性能校车预约平台，旨在解决学生及教职工的日常通勤预约需求。系统采用了 **简洁架构 (Clean Architecture)** 设计，确保代码的可维护性与可扩展性，同时针对高并发场景进行了专项优化。

## 🚀 核心技术栈

- **语言**: Go (1.21+)
- **Web 框架**: Gin
- **ORM**: Gorm (MySQL 8.0)
- **缓存**: Redis 7.0 (Lua 脚本)
- **消息队列**: RabbitMQ (异步下单)
- **日志**: slog (结构化日志)
- **认证**: JWT (JSON Web Token)
- **并发工具**: Singleflight

## 🏗️ 架构设计

系统遵循简洁架构原则，划分为以下层次：

1.  **Domain (领域层)**: 核心业务实体 (User, Bus, Order) 及存储接口定义。
2.  **Usecase (用例层)**: 核心业务逻辑，如预约校验、库存管理、状态流转。
3.  **Delivery (交付层)**: 处理外部输入，包括 HTTP API (Gin) 和 MQ 消费者 (RabbitMQ Consumer)。
4.  **Repository (基础设施层)**: 具体的数据库实现 (Gorm) 和缓存实现 (Redis)。

### 目录结构

```text
.
├── cmd/server/             # 程序入口，负责依赖注入与服务启动
├── internal/
│   ├── domain/             # 领域实体与接口定义
│   ├── usecase/            # 业务逻辑实现 (Business Logic)
│   ├── delivery/
│   │   ├── http/           # HTTP Handlers & Middleware
│   │   └── mq/             # RabbitMQ 消费者 (Order Processing)
│   └── repository/         # 数据持久化实现 (MySQL, Redis)
├── pkg/                    # 基础设施工具 (Database, MQ)
├── config/                 # 配置管理
├── scripts/                # SQL 初始化脚本、Lua 脚本
└── docker-compose.yaml     # 环境一键启动
```

## 🔥 核心业务优化方案

### 1. 秒杀级库存扣减 (Redis + Lua)
在高并发预约场景下，利用 Redis 的单线程原子性，通过 Lua 脚本实现“查询库存-判断-扣减”的原子操作，从根源上杜绝超卖问题。

### 2. 缓存击穿防护 (Singleflight)
针对热门班次的频繁并发查询，引入 `golang.org/x/sync/singleflight`。在缓存失效时，确保同一时间内只有一个请求去查询数据库，并将结果共享给其他并发请求，极大地降低了数据库压力。

### 3. 异步解耦与削峰填谷 (RabbitMQ)
- **极速响应**: 用户下单后，系统在 Redis 扣减成功后立即通过 MQ 发送任务并返回“处理中”，不阻塞 HTTP 连接。
- **可靠持久化**: 消费者通过可靠消费模式将订单持久化到 MySQL，并更新最终库存状态。

### 4. 全链路控制与优雅退出
利用 Go 原生 `Context` 实现超时控制，结合信号监听实现服务的 **Graceful Shutdown**，确保在服务重启时正在处理的消息能够被正确 Ack 或持久化。

## 🛠️ 快速开始

### 1. 启动基础设施
使用 Docker Compose 一键启动 MySQL, Redis, RabbitMQ：
```bash
docker-compose up -d
```

### 2. 配置文件
修改 `config.yaml`（若无则从 `config.yaml.example` 复制）：
```bash
cp config.yaml.example config.yaml
```

### 3. 运行服务
```bash
go run cmd/server/main.go
```

## 接口预览

| 方法 | 路径 | 说明 | 认证 |
| :--- | :--- | :--- | :--- |
| POST | `/api/v1/users/register` | 用户注册 | 否 |
| POST | `/api/v1/users/login` | 用户登录 (返回 JWT) | 否 |
| GET | `/api/v1/users/profile` | 获取个人信息 | 是 |
| GET | `/api/v1/buses` | 获取所有班次列表 | 否 |
| GET | `/api/v1/buses/:id` | 获取班次详情 | 否 |
| POST | `/api/v1/orders` | 提交预约订单 (异步) | 是 |
| GET | `/api/v1/orders` | 查询用户订单列表 | 是 |
| POST | `/api/v1/orders/:id/cancel` | 取消订单 | 是 |

## 🧪 压测说明
建议使用 **JMeter** 进行压力测试。在 500 QPS 下，通过 Redis + Lua 方案可实现“零超卖”且平均响应延迟控制在 50ms 以内。
