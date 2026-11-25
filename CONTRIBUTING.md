# Contributing

Every change to ripple goes through a pull request, and the owner reviews and merges every pull request. This file describes the workflow. The rules for the code itself are in [docs/conventions.md](docs/conventions.md).

## Branches

- `main` is always deployable. Nobody pushes to it directly.
- Work happens on short-lived branches that are merged within a few days.
- Branch names start with the change type: `feat/home-timeline-pagination`, `fix/duplicate-like`, `docs/adr-object-storage`.

## Commits

Commit messages follow [Conventional Commits](https://www.conventionalcommits.org/):

```text
<type>(<scope>): <summary>

<body: what changed and why>
```

- **Types:** `feat`, `fix`, `perf`, `refactor`, `test`, `docs`, `build`, `ci`, `chore`.
- **Scopes:** `content`, `feed`, `media`, `deploy`, `loadtest`, `docs`, `ci`.
- **Summary:** imperative mood, lowercase, no period, at most 72 characters. It should complete the sentence "If applied, this commit will ...".
- **Body:** explain why the change was made, not only what changed. Wrap lines at 72 characters.
- **Breaking changes:** add `!` after the scope and a `BREAKING CHANGE:` footer.

Examples:

```text
feat(feed): add cursor pagination to home timeline
fix(content): keep one like per account under concurrent requests
perf(feed): batch timeline writes during fan-out
docs(adr): choose object storage for media
```

## Pull requests

- One pull request makes one logical change. Aim for fewer than 400 changed lines, not counting generated files.
- The title is a Conventional Commit summary. Pull requests are squash-merged, so the title becomes the commit on `main`.
- The description follows the [pull request template](.github/pull_request_template.md), with every section filled in.
- CI must pass before merging.
- The owner approves and merges. The branch is deleted after merging.

## Definition of done

A change is done when:

- Tests cover the new behavior, including failure cases and, where relevant, the [charter invariants](docs/charter.md#invariants).
- CI passes.
- Logs and metrics cover the new behavior.
- Documentation and ADRs reflect any change in behavior or decisions.
- User-facing changes are deployed to production through the pipeline.

## Review checklist

Reviewers check that the change:

- Stays correct under concurrent requests, retries, and redelivered events
- Handles failures: timeouts, retries with backoff, and what happens when a dependency is down
- Keeps every charter invariant
- Validates input and checks authorization
- Has tests for the failure paths, not only the happy path
- Uses clear names and clear English
