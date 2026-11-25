# CLAUDE.md

Guidance for Claude when working in this repository.

## Read first

- [docs/charter.md](docs/charter.md): scope, invariants, and constraints
- [docs/conventions.md](docs/conventions.md): naming, APIs, data, events, observability, and testing
- [CONTRIBUTING.md](CONTRIBUTING.md): branches, commits, and pull requests
- [docs/adr/](docs/adr/README.md): decisions already made. Never contradict an accepted ADR; propose a new one instead.

## Division of work

The owner uses this project to build backend skills, so the core logic is theirs to write.

- **The owner writes** the core logic and its tests: fan-out, timeline pagination, counters, idempotency handling, and the media processing pipeline.
- **Claude writes** scaffolding, build files, CI/CD workflows, Dockerfiles and Compose files, Nginx and host configuration, observability wiring, load-test tooling, documentation, and every pull request description.
- **Claude does not implement owner-owned logic** unless the owner explicitly asks. It explains, gives hints, suggests test cases, and reviews instead.
- **Claude reviews the owner's code** against the [review checklist](CONTRIBUTING.md#review-checklist). It also points out unnatural English in commit messages, comments, and identifiers, and suggests a more natural version. The owner decides whether to change it.

## Writing

- Write everything in the repository in English, using the preferred terms in [docs/conventions.md](docs/conventions.md#language).
- Keep sentences short and plain.

## Pull requests

- Follow [.github/pull_request_template.md](.github/pull_request_template.md) and fill in every section.
- Write the title as a Conventional Commit summary, since it becomes the squash commit on `main`.
- Include load-test numbers when a change affects performance.
