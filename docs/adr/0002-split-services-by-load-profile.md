# 0002. Split services by load profile

- Status: Accepted
- Date: 2026-11-25

## Context

A feed system does three kinds of work, and each has a very different load profile:

- **Writes and domain rules:** creating posts, following, liking, and sending notifications. Moderate volume, transactional, and full of business rules.
- **Timeline reads and fan-out:** reads outnumber writes by orders of magnitude, and a post from an account with many followers can mean a write to every follower's timeline.
- **Image processing:** CPU-heavy and bursty. It must not slow down request handling.

The project also wants to use Java and Go where each fits best, and to practice the problems that only appear between services: timeouts, retries, and keeping data consistent across service boundaries.

## Decision

The system has three components:

| Component | Language | Responsibility |
|---|---|---|
| `content` | Java 17, Spring Boot 4 | Accounts, posts, follows, likes, notifications, media metadata, and upload authorization. It is the system of record, stored in MySQL. |
| `feed` | Go 1.27 | Home timelines and the hot list, stored in Redis. It serves the timeline API and runs the fan-out workers. |
| `media` | Go 1.27 | Image processing. It has no API apart from health checks and metrics. |

`content` publishes domain events to Kafka through a transactional outbox. `feed` and `media` consume those events, and `media` reports its results back through events. No component reads another component's data store.

## Alternatives considered

- **A Java modular monolith, with Go workers for fan-out and media.** Fewer moving parts, but it leaves almost no cross-service problems to practice, and gives Go a small role.
- **Fine-grained microservices**, with separate services for accounts, posts, the follow graph, feeds, likes, notifications, and media. It resembles a large company's architecture, but for a single developer it is mostly operational overhead.

## Consequences

- Components scale and fail independently. If `media` falls behind, text posts and timeline reads keep working.
- Data crosses component boundaries only through events and APIs, so consistency between components is eventual. This needs deliberate design: idempotent consumers, reconciliation jobs, and read-your-writes handling, as the [charter's invariants](../charter.md#invariants) require.
- Debugging a request that crosses components needs distributed tracing.
- How `feed` turns post IDs into full posts for a timeline (calling `content`, or keeping its own copy built from events) needs its own ADR.
- Three deployables mean three build pipelines, more dashboards, and more deploy steps. Adding a fourth component needs a new ADR.
