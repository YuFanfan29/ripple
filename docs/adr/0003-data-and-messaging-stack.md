# 0003. Use MySQL, Redis, and Kafka

- Status: Accepted
- Date: 2026-11-25

## Context

ripple needs a durable system of record, a fast store for timelines and counters, a stream of events that several consumers can read at their own pace, and storage for images. The project should use what large internet companies commonly run, so that it meets the problems those companies deal with.

## Decision

- **MySQL 8.4 LTS**, with one primary and one replica, is the system of record. Reading from the replica means the code must handle replication lag, for example to keep read-your-writes.
- **Redis 8** stores home timelines, counters, caches, and rate limits. Redis is never the only copy of anything: its data can be lost and rebuilt.
- **Kafka 4 in KRaft mode** carries domain events between components.
- **Self-hosted, S3-compatible object storage** holds images. The product, Garage or SeaweedFS, is chosen in a separate ADR before media work starts. MinIO is excluded: its community edition was archived in 2026, and its images were removed from Docker Hub.

## Alternatives considered

- **PostgreSQL instead of MySQL.** An excellent database, but MySQL with InnoDB is more common at large internet companies, and the project wants to work with InnoDB specifics: clustered indexes, MVCC, gap locks, and binlog-based replication.
- **RocketMQ instead of Kafka.** It has built-in delayed and transactional messages, but Kafka is more widely used across the industry, and the transactional outbox covers reliable publishing.
- **Valkey instead of Redis.** An open-source fork that is compatible with Redis. Either works for a self-hosted project, and switching later needs no code changes.

## Consequences

- MySQL replication, Kafka, and Redis together need careful memory planning on small hosts. Every container gets an explicit memory limit.
- Every Redis structure needs a documented way to rebuild it, through the owning service's source data, another service's API, or by replaying events.
- All events leave `content` through the outbox relay, so the relay is critical. Its lag is monitored and alerted on.
