# ripple charter

This document defines what ripple is, what it must guarantee, and the constraints it is built under. Change it through a pull request, and record the reason for any significant change in an [ADR](adr/README.md).

## Purpose

ripple is a learning and portfolio project. It exists to give its owner hands-on experience with the backend problems that content platforms face, and to produce evidence of that experience: measured results, failure drills, and documented decisions.

The product is small on purpose. A feature belongs in scope only if it brings an engineering problem worth solving.

## Scope

### In scope

- **Accounts.** Sign-up and sign-in. Kept simple; not a focus of the project.
- **Posts.** Text with up to 4 images. Accounts can create, delete, and view posts. Each account has a user timeline of its own posts.
- **Media.** Clients upload images directly to object storage using presigned URLs. The `media` worker validates each image, strips its metadata, fixes its orientation, and creates resized variants. An image becomes visible only after processing succeeds. Uploads that are not attached to a post within 24 hours are deleted.
- **Follows.** Follow and unfollow, follower and following counts, and paginated follower and following lists.
- **Home timeline.** Posts from followed accounts, newest first, with cursor pagination. Fan-out is hybrid: fan-out on write for most accounts, and fan-out on read for celebrity accounts.
- **Likes.** Like and unlike a post, and show like counts.
- **Notifications.** New followers and likes, with likes grouped ("A and 12 others liked your post"), and an unread count.
- **Hot list.** Popular recent posts, ranked by a score that decays over time.
- **Demo page.** A minimal web page that shows the system working. It is not a product frontend.

### Later

Comments, search, real-time updates, and video. Each one needs an ADR before work starts.

### Out of scope

- Recommendations or machine-learned ranking
- Direct messages
- A full web frontend or mobile apps
- Multi-region deployment
- Kubernetes, unless a later ADR shows a concrete need

## Invariants

These rules must always hold, including under load and during failures. Each one is covered by automated tests, and is checked again after every load test and failure drill.

1. **One like per account per post.** Repeated or concurrent like requests leave exactly one like. Unlike is idempotent.
2. **Counts converge.** Like and follower counts may lag behind the source of truth, but must match it within 5 seconds. A reconciliation job finds and repairs drift, and lag above 5 seconds raises an alert.
3. **Deleted posts are never served.** After a delete request succeeds, no API returns the post.
4. **Read-your-writes.** Authors see their own new posts in their home timeline immediately, even though fan-out is asynchronous.
5. **Stable pagination.** Paging through a timeline never returns a post twice and never skips one, even while new posts arrive.
6. **Safe media.** No image is served before it has been validated (type detected from its content, limits on file size and pixel count) and stripped of location and other EXIF metadata. An account can attach only media it uploaded.
7. **No duplicate side effects.** Every event consumer is idempotent. A redelivered event never creates a duplicate notification or counts anything twice.

## Operational constraints

- **Pipeline-only delivery.** Every change reaches production through the CI/CD pipeline. There are no manual deploys.
- **Zero-downtime deploys.** A deploy causes no failed requests. A k6 load test runs during deploys to prove it.
- **Automatic rollback.** A release that fails its health checks is rolled back automatically.
- **One-command setup.** `make up` starts the whole system from a clean clone, and `make ci` runs the same checks as CI.
- **Infrastructure as code.** Production runs on one or two small Linux VPS hosts. All host configuration lives in this repository, and any host can be rebuilt from scratch.
- **Own domain.** The system is served over TLS on its own domain name, and images can be served through a CDN.
- **Realistic load tests.** Load tests run from a separate machine, against a generated dataset with a skewed follower distribution: most accounts have few followers, and a few have very many.
- **Measured performance.** Performance targets are set after the first baseline load test. From then on, `docs/performance.md` records the targets, the results, and every bottleneck found and fixed.
- **Multi-architecture images.** Container images are built for `linux/amd64` and `linux/arm64`.

## Technical constraints

| Area | Choice |
|---|---|
| `content` service | Java 17 (LTS), Spring Boot 4, Maven Wrapper |
| `feed` service, `media` worker | Go 1.27 |
| System of record | MySQL 8.4 LTS, one primary and one replica |
| Timelines, counters, caches | Redis 8 |
| Events | Kafka 4 in KRaft mode |
| Object storage | Self-hosted and S3-compatible: Garage or SeaweedFS, chosen in an ADR before media work starts |
| Edge | Nginx |
| Deployment | Docker Compose on Linux VPS hosts |
| CI/CD | GitHub Actions |
| Observability | Prometheus, Grafana, OpenTelemetry |
| Load testing | k6 |

Every component runs in a container. Services use object storage only through the S3 API, so the storage product can change without code changes. [ADR 0002](adr/0002-split-services-by-load-profile.md) explains the service split, and [ADR 0003](adr/0003-data-and-messaging-stack.md) explains the data and messaging choices.

## Open decisions

- GitHub owner, which sets the Go module path and the Java group ID
- License
- Domain name
- Object storage product (Garage or SeaweedFS)
