# Architecture decision records

Each record describes one significant decision: its context, the decision, the alternatives, and the consequences. A decision is significant if it is expensive to reverse, affects more than one component, or chooses between real alternatives.

To add a record, copy [0000-template.md](0000-template.md) and use the next number. Accepted records are not rewritten. When a decision changes, a new record supersedes the old one, and only the old record's status line changes.

| No. | Decision | Status |
|---|---|---|
| [0001](0001-record-architecture-decisions.md) | Record architecture decisions | Accepted |
| [0002](0002-split-services-by-load-profile.md) | Split services by load profile | Accepted |
| [0003](0003-data-and-messaging-stack.md) | Use MySQL, Redis, and Kafka | Accepted |
