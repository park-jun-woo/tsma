# tsma

Regression defense for legacy code. Indexes every function, detects missing tests, measures branch coverage, and guides LLM agents to fill the gaps — one function at a time.

## What tsma does

1. **Shows where tests are missing.** Indexes every function, matches against test files by convention, and tells you exactly which functions have no tests.
2. **Uncovered branch feedback makes LLM tests dramatically better.** Without feedback, LLM writes 60-70% coverage. With specific line numbers ("line 41, 44, 70 uncovered"), LLM reaches 100% in one shot.
3. **Single-command loop for LLM agents.** `tsma next` handles detection, test execution, coverage measurement, and progress tracking. The agent just repeats one command until done.

Validated on a real project (527 functions): 246 reached 100% (PASS), 281 accepted at best-effort (DONE), 0 remaining (TODO).

## Install

```bash
make install
```

## How it works

`tsma next` is the only command you need. It drives the entire loop:

```
$ tsma next          # shows the next function without a test
  → write a test
$ tsma next          # detects the new test, runs it, measures coverage
  → 100%? PASS, moves to next function
  → <100%? shows uncovered branches, gives you one more shot
$ tsma next          # re-measures after your fix
  → improved or not, marks DONE, moves on
```

Repeat until `All functions complete!`.

## Commands

```
tsma next              the entire workflow (detect → test → coverage → advance)
tsma list [--page N]   all functions with status
tsma status            progress summary
tsma reset --all       delete session
```

## Status

| Status | Meaning |
|---|---|
| **TODO** | no test file, or coverage < 100% and due this round |
| **DONE** | not 100% but best effort this round |
| **PASS** | 100% branch coverage |

## Example session

```
$ tsma next

No session found. Analyzing project...
Detected: go
Found 527 functions
Session created.

Login  TODO
  file: internal/api/auth/handler.go:86-119
  test: internal/api/auth/handler_test.go (not found)
```

After writing a test:

```
$ tsma next

Login  testing...
  [1/2] go test: PASS
  [2/2] coverage: 100% (9/9)
PASS ✓

GetProfile  TODO
  file: internal/api/auth/handler.go:121-148
  test: (not found)
```

When coverage is not 100%:

```
$ tsma next

ListContracts  testing...
  [1/2] go test: PASS
  [2/2] coverage: 65% (11/17)
  UNCOVERED:
    line 41 — if params.Status != nil
    line 44 — if params.BuildingId != nil
    line 70 — if err != nil (CountSummary)
  ▶ Cover the uncovered branches. Next tsma next will re-measure.
```

## Principles

1. **Convention-based test matching.** Go: `handler.go` → `handler_test.go` in the same directory. TS: `.test.ts` / `.spec.ts`. Python: `test_` prefix.
2. **Session is cache, source files are truth.** If a test file is deleted, the function reverts to TODO regardless of what session.json says.
3. **Generated code is excluded.** `*_gen.go`, `*.pb.go` are not indexed.
4. **Coverage is measured, not enforced.** 100% is the goal, but unreachable branches (no DI, external dependencies) are accepted as DONE after one retry.

## Why some functions can't reach 100%

Whether a function can reach 100% branch coverage depends on how it receives its dependencies.

**Interface (mockable) → 100% achievable:**

```go
type Handler struct {
    svc AuthSvc              // interface — can be replaced with a mock
}
```

In tests, you inject a mock that returns whatever you need:

```go
svc := mocks.NewMockAuthSvc(ctrl)
svc.EXPECT().Login(...).Return(result, nil)   // success path
svc.EXPECT().Login(...).Return(nil, err)      // error path
```

Every branch is reachable because you control all inputs and outputs.

**Concrete type (not mockable) → stuck below 100%:**

```go
type Handler struct {
    svc *service.SMSImportService    // struct pointer — cannot be replaced
}
```

The real implementation runs with all its internal dependencies (DB, external APIs). You cannot make it return a specific error or a specific result. Branches that depend on those outcomes are unreachable in unit tests.

**tsma's response:** After one retry with uncovered branch feedback, tsma marks the function as DONE with the achieved coverage. This is not a tool limitation — it reflects the code's testability. Introducing an interface (DI) would make the function fully testable, but that requires modifying the source code.

## Language support

| Language | Indexer | Test runner | Coverage |
|---|---|---|---|
| Go | `go/ast` | `go test` | `go test -coverprofile` |
| TypeScript | regex | `npx vitest` / `npx jest` | `c8` / `istanbul` |
| Python | regex | `pytest` | `coverage.py` |

## For LLM agents

Give this instruction to your agent:

```
1. Run `tsma next`
2. If TODO — read the function, write a test
3. If test fails — read the error, fix the test
4. If uncovered branches shown — add tests for those branches
5. If PASS/DONE — next function is automatically shown
6. Repeat until "All functions complete!"
```

The agent only needs to know one command: `tsma next`.

## Reports

- [REPORT.md](REPORT.md) — tsma self-test (115 functions, PASS 87, DONE 28, coverage 95%)
- [REPORT.juicer.md](REPORT.juicer.md) — juicer project test (140 functions, PASS 114, DONE 26, coverage 96%)

## Related work

[Diffblue Cover](https://www.diffblue.com/) uses reinforcement learning + symbolic execution to generate Java unit tests in a similar generate → verify → feedback loop. tsma takes the same core insight — deterministic verification driving iterative generation — but uses LLMs as the generator, works across Go/TypeScript/Python, and is open source.
