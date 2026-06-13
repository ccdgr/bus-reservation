# 校车预定平台 (Bus Reservation Platform)

本项目是一个基于 Go 语言和 React 开发的高性能校车预约平台，旨在解决学生及教职工的日常通勤预约需求。系统采用了 **简洁架构 (Clean Architecture)** 设计，后端针对高并发场景进行了深度优化，前端适配移动端 WebView 嵌入。

## 🚀 核心技术栈

### 后端 (Backend)
- **核心语言**: Go (1.21+)
- **Web 框架**: Gin (高性能路由)
- **数据库**: MySQL 8.0 (Gorm ORM)
- **缓存**: Redis 7.0 (Lua 脚本实现原子性库存扣减)
- **消息队列**: RabbitMQ (异步下单、延时取消、削峰填谷)
- **并发优化**: Singleflight (缓存击穿防护), Context (全链路超时控制)
- **日志系统**: slog (结构化 JSON 日志)

### 前端 (Frontend)
- **框架**: React 19 + TypeScript
- **UI 库**: MUI (Material UI v6) - 移动端优先设计
- **路由**: React Router v7
- **状态管理**: React Context (Auth 认证流)
- **通信**: Axios (带 JWT 自动注入与 401 拦截)

## 🏗️ 架构设计

系统遵循简洁架构原则，划分为以下层次：

1.  **Domain (领域层)**: 核心业务实体 (User, Bus, Order) 及状态机逻辑。
2.  **Usecase (用例层)**: 核心业务流程实现（如：预约、支付、核验、自动过期）。
3.  **Delivery (交付层)**: HTTP API (Gin Handlers) 与消息消费者 (MQ Consumers)。
4.  **Repository (基础设施层)**: 持久化实现 (Gorm, Redis, MQ Producers)。

### 目录结构

```text
.
├── cmd/server/             # 后端入口 (依赖注入与启动)
├── internal/
│   ├── domain/             # 领域实体、状态机、接口定义
│   ├── usecase/            # 业务逻辑实现
│   ├── delivery/
│   │   ├── http/           # RESTful API、Auth 中间件
│   │   └── mq/             # RabbitMQ 消费者 (异步入库、延时取消)
│   └── repository/         # 数据持久化实现 (MySQL, Redis)
├── frontend/               # 前端 React 项目
│   └── src/
│       ├── pages/          # 首页搜索、结果页、订单管理、个人中心
│       └── context/         # 全局认证状态
├── pkg/                    # 基础设施工具 (Database, MQ)
├── scripts/                # SQL 初始化 (已针对 utf8mb4 优化)、Lua 脚本
└── docker-compose.yaml     # 基础设施一键启动
```

## 🔥 核心业务优化方案

### 1. 秒杀级库存扣减 (Redis + Lua)
利用 Redis 的单线程原子性，通过 Lua 脚本实现“查询库存-判断-扣减”的闭环操作，确保在高并发压力下实现 **“零超卖”**。

### 2. 缓存击穿防护 (Singleflight)
针对热门班次的并发查询，引入 `singleflight`。在缓存失效时，合并重复的数据库读请求，确保数据库压力不因并发激增而崩溃。

### 3. 异步解耦与延时取消 (RabbitMQ)
- **异步下单**: 用户点击预定后，系统在 Redis 扣减成功后立即通过 MQ 发送任务并返回，大幅提升响应速度。
- **延时取消**: 利用 **Dead Letter Exchange (DLX)** 实现延时队列。订单创建后自动进入 15 分钟倒计时，若超时未支付，系统将自动取消订单并归还库存。

### 4. 严谨的订单状态机
定义了 `待支付`、`待核验`、`已取消`、`已过期`、`已核验` 五大状态。通过状态机校验确保业务流转的安全性和一致性。

## 🛠️ 快速开始

### 1. 环境准备
使用 Docker Compose 一键启动 MySQL, Redis, RabbitMQ：
```bash
docker-compose up -d
```
*注意：若首次运行出现乱码，请执行 `docker-compose down -v` 后重新启动。*

### 2. 启动后端
```bash
cp config.yaml.example config.yaml
# 修改 config.yaml 中的 DSN 和密钥配置
go run cmd/server/main.go
```

### 3. 启动前端
```bash
cd frontend
npm install
npm run dev -- --host
```

## 🎥 功能特性预览

- **智能搜索**: 支持按出发地、目的地、日期（禁止选择过去日期）精确筛选班次。
- **全流程支付**: 模拟一卡通、微信、支付宝多种支付方式，体验真实的支付反馈。
- **乘车核验**: 生成核验二维码，模拟扫码上车后的状态实时同步。
- **订单追踪**: 实时查看订单状态，支持待核验订单的退票处理。

## 🧪 压测性能
建议使用 **JMeter** 进行测试。在单机环境下，核心下单接口可支持 1000+ QPS，平均延迟低于 50ms。
