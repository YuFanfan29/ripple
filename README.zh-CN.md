# ripple

[English](README.md) | 简体中文

ripple 是一个内容 feed 系统。用户可以发布文字和图片、互相关注、给帖子点赞，并在首页时间线（home timeline）上浏览所关注账号的帖子。

产品刻意做得很小，工程上的挑战却不小。ripple 要逐一解决真实 feed 系统会遇到的问题：超大粉丝量账号的 fan-out（扇出）、读多写少的流量、高频写入下的计数器、热点内容，以及异步图片处理。整个系统运行在 Linux 主机上，并通过 CI/CD 流水线发布。

## 当前状态

脚手架阶段。`content`、`feed` 和 `media` 三个服务都能启动、响应健康检查，并能优雅关闭（graceful shutdown）。目前还没有业务逻辑。

## 快速开始

需要 Go 1.27 和 JDK 17。Maven 通过 Maven Wrapper 随仓库提供，不需要单独安装。本地可以不装 [golangci-lint](https://golangci-lint.run/)，CI 总会运行它。

```bash
make test   # 运行所有 Go 和 Java 测试
make lint   # 检查代码格式并运行静态分析
make ci     # 运行与 CI 相同的检查
```

启动一个服务并检查健康状态：

```bash
go run ./services/feed                         # 监听 :8081
(cd services/content && ./mvnw spring-boot:run) # 监听 :8080
curl localhost:8081/readyz
```

## 架构

| 组件 | 语言 | 职责 |
|---|---|---|
| `content` | Java 17、Spring Boot 4 | 账号、帖子、关注、点赞、通知和媒体元数据。整个系统的权威数据源（system of record）。 |
| `feed` | Go 1.27 | 首页时间线、fan-out 和热榜。 |
| `media` | Go 1.27 | 图片校验、去除元数据和调整尺寸。 |

- **基础设施：** MySQL 8.4（一主一从）、Redis 8、Kafka 4、兼容 S3 的对象存储，以及 Nginx。
- **交付：** GitHub Actions 用 Docker Compose 部署到 Linux VPS 主机。
- **可观测性：** Prometheus、Grafana 和 OpenTelemetry。
- **压测：** k6。

[ADR 0002](docs/adr/0002-split-services-by-load-profile.md) 解释了系统为什么这样拆分。

## 文档

以下文档都是英文：

- [Charter](docs/charter.md)（项目章程）：目标、范围、不变量和约束
- [Conventions](docs/conventions.md)（开发规范）：命名、API、数据、事件、可观测性和测试
- [Architecture decision records](docs/adr/README.md)（架构决策记录）：重要的决策及其原因
- [Contributing](CONTRIBUTING.md)（贡献指南）：分支、提交、pull request 和代码评审

---

本文是 [README.md](README.md) 的中文翻译。两者不一致时，以英文版为准。
