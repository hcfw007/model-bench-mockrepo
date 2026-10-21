# model-bench-mockrepo

[中文](README.md) | **English**

The testbed for [model-bench](https://github.com/hcfw007/model-bench).

A de-identified multi-service repository that reproduces the shape of my real work:
cross-service localization, minimal edits, evidence gathering, SOP execution, and
subagent collaboration. None of it is real business code.

## Why a separate repository

The main repo uses `git worktree` to isolate every case, so the testbed needs its own
git history. One case requires finding the commit that introduced a problem purely
from `git log`, so the history intentionally keeps a defective change.

## Layout

```
go-projects/
  some-service/          task creation, confirmation, polling, tag marking
    internal/biz/bulk_send_v1/   business layer in a Kratos-style split
  auth-service/          permission groups, member groups, NEW badges
  contact-data-service/  contact records and tags
ts-projects/
  api-bff/               BFF: material validation, confirm-item projection
  web-console/           frontend pages
.tasks/bulk-send-v1/     cross-session task archive (STATUS / record / case tree)
fixtures/                synthetic logs, mongo documents, a hanging script
```

## What is seeded

Each service contains a handful of defects isomorphic to real ones:

- completion condition: does a job converge when some targets failed
- auto tag marking: is the stored config ever executed
- confirm endpoint: is a canceled job's error distinguishable
- frontend list: does it consume the confirmability field the BFF already computes
- permissions and badges: missing registration, truncated pagination
- task archive: two records contradict each other (tests whether the agent stops and asks)

Exact assertions live in the main repo's `cases/`; they are not repeated here.

## Verify

```bash
go build ./...
go test ./... -count=1
```

Zero external dependencies on the Go side.

## Constraints

- no real remote configured, `git push` is forbidden
- no calls to real messaging or broadcast endpoints
- all IDs, logs, and customer names are synthetic
