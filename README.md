# ripple

English | [简体中文](README.zh-CN.md)

ripple is a content feed system. People post text and photos, follow each other, like posts, and read a home timeline made from the accounts they follow.

The product is small on purpose; the engineering is not. ripple is built to work through the problems that real feed systems face: fan-out for accounts with very many followers, read-heavy traffic, counters under heavy writes, hot content, and asynchronous image processing. All of it runs on Linux hosts and ships through a CI/CD pipeline.

## Status

Scaffolding. `content`, `feed`, and `media` start, serve health checks, and shut down gracefully. There is no business logic yet.

## Quickstart

You need Go 1.27 and JDK 17. Maven comes with the repository through the Maven Wrapper. [golangci-lint](https://golangci-lint.run/) is optional locally; CI always runs it.

```bash
make test   # run all Go and Java tests
make lint   # check formatting and run static analysis
make ci     # run the same checks as CI
```

To run a service and check its health:

```bash
go run ./services/feed                         # serves on :8081
(cd services/content && ./mvnw spring-boot:run) # serves on :8080
curl localhost:8081/readyz
```

## Architecture

| Component | Language | Responsibility |
|---|---|---|
| `content` | Java 17, Spring Boot 4 | Accounts, posts, follows, likes, notifications, and media metadata. The system of record. |
| `feed` | Go 1.27 | Home timelines, fan-out, and the hot list. |
| `media` | Go 1.27 | Image validation, metadata stripping, and resizing. |

- **Infrastructure:** MySQL 8.4 (primary and replica), Redis 8, Kafka 4, S3-compatible object storage, and Nginx.
- **Delivery:** GitHub Actions deploys with Docker Compose to Linux VPS hosts.
- **Observability:** Prometheus, Grafana, and OpenTelemetry.
- **Load testing:** k6.

[ADR 0002](docs/adr/0002-split-services-by-load-profile.md) explains why the system is split this way.

## Documentation

- [Charter](docs/charter.md): goals, scope, invariants, and constraints
- [Conventions](docs/conventions.md): naming, APIs, data, events, observability, and testing
- [Architecture decision records](docs/adr/README.md): significant decisions and their reasons
- [Contributing](CONTRIBUTING.md): branches, commits, pull requests, and reviews
