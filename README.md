# 西工大校车预定平台 (NWPU Bus Reservation Platform)

本项目是一个基于 Go 语言开发的高性能校车预约平台，旨在解决学生及教职工的日常通勤预约需求。系统采用了简洁架构（Clean Architecture）设计，确保代码的可维护性与可扩展性，同时针对高并发场景进行了专项优化。

## 核心技术栈

- **后端**: Go, Gin (Web 框架)
- **数据库**: MySQL (Gorm), Redis
- **消息队列**: RabbitMQ
- **并发优化**: Singleflight, Redis + Lua 脚本
- **测试/压测**: JMeter
- **前端**: Vue

## 架构设计

系统遵循简洁架构原则，划分为以下层次：

1. **Domain (领域层)**: 定义核心业务实体（Entity）及持久化接口（Repository Interface）。
2. **Usecase (用例层)**: 封装业务逻辑，处理具体业务流程（如预约校验、订单流转）。
3. **Delivery/Interface (接口层)**: 处理外部输入，包括 HTTP 请求（Gin）及消息队列消费。
4. **Repository/Infrastructure (基础设施层)**: 实现数据持久化逻辑（Gorm, Redis, MQ）。

### 目录结构

```text
.
├── cmd/
│   └── server/             # 程序入口，负责依赖注入与服务启动
├── internal/
│   ├── domain/             # 领域实体与接口定义
│   ├── usecase/            # 业务逻辑实现
│   ├── delivery/           # HTTP 接口、MQ 消费者
│   │   ├── http/           # Gin Handlers
│   │   └── mq/             # RabbitMQ Consumers
│   └── repository/         # 数据持久化实现 (MySQL, Redis)
├── pkg/                    # 通用工具类 (Logger, Response, etc.)
├── config/                 # 配置文件及加载逻辑
├── api/                    # 接口文档 (Swagger)
├── scripts/                # Lua 脚本、数据库迁移脚本
└── README.md
```

## 核心业务优化方案

### 1. 高并发库存扣减 (Redis + Lua)
利用 Redis 的单线程原子性，通过 Lua 脚本实现“查询库存-判断-扣减”的原子操作，从根源上杜绝超卖问题。

### 2. 缓存击穿防护 (Singleflight)
针对热门班次的并发查询，使用 `golang.org/x/sync/singleflight` 机制，确保同一时刻只有一个请求去查询数据库，减轻后端压力。

### 3. 异步解耦与可靠投递 (RabbitMQ)
- **异步处理**: 下单成功后通过 MQ 异步通知其他系统或处理后续逻辑（如超时自动取消）。
- **可靠性**: 采用生产者确认模式（Publisher Confirms）及消费者幂等性设计，确保消息不丢失且不重复消费。

### 4. 链路控制与优雅退出
利用 Go 原生 `Context` 实现全链路超时控制，并结合信号监听实现服务的优雅退出（Graceful Shutdown），确保正在处理的订单不中断。

## 快速开始

(待完善...)
