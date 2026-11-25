# 0001. Record architecture decisions

- Status: Accepted
- Date: 2026-11-25

## Context

ripple is built by one person over several months. Early decisions, such as how the services are split and which data stores to use, shape everything that follows, and their reasons are easy to forget. Being able to explain those reasons is also part of the project's purpose (see the [charter](../charter.md#purpose)).

## Decision

Every significant decision is recorded as an architecture decision record (ADR) in `docs/adr/`, using the format in [0000-template.md](0000-template.md). A decision is significant if it is expensive to reverse, affects more than one component, or chooses between real alternatives.

An accepted ADR is not rewritten. When a decision changes, a new ADR supersedes the old one, and the old one's status line points to the new one.

## Alternatives considered

- **No written record.** Fastest at first, but the reasons get lost, and settled questions get argued again.
- **One long design document.** Harder to keep current, and it hides when and why each decision was made.

## Consequences

- Every significant change includes an ADR or links to one.
- The list of ADRs doubles as a history of how the system evolved.
