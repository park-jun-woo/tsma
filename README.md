# tsma

[![Version](https://img.shields.io/badge/version-v0.5.0-blue.svg)](https://github.com/park-jun-woo/tsma/releases)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![skills.sh](https://skills.sh/b/park-jun-woo/tsma)](https://skills.sh/park-jun-woo/tsma)

Regression defense for legacy code. Indexes every function, detects missing tests, measures branch coverage, and guides LLM agents to fill the gaps — one function at a time.

## What tsma does

1. **Shows where tests are missing.** Indexes every function, matches against test files by convention, and tells you exactly which functions have no tests.
2. **Uncovered branch feedback makes LLM tests dramatically better.** Without feedback, LLM writes 60-70% coverage. With specific line numbers ("line 41, 44, 70 uncovered"), LLM reaches 100% in one shot.
3. **Single-command loop for LLM agents.** `tsma next` handles detection, test execution, coverage measurement, and progress tracking. The agent just repeats one command until done.

Validated on a real project (527 functions): 246 reached 100% (PASS), 281 accepted at best-effort (DONE), 0 remaining (TODO).

## Quick Start

```bash
npx skills add park-jun-woo/tsma
```

## Install

```bash
go install github.com/park-jun-woo/tsma/cmd/tsma@latest
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
tsma rescan            re-sync the function set with current source (keeps progress)
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
  ▶ Cover the uncovered branches.
  ▶ After completing the test, run `tsma next`.
```

## Principles

1. **Convention-based test matching.** Go: `handler.go` → `handler_test.go` in the same directory. TS: `.test.ts` / `.spec.ts`. Python: `test_` prefix. Rust: in-file `#[cfg(test)] mod tests` or `tests/*.rs`. Java: `src/main/java/…/Foo.java` → `src/test/java/…/FooTest.java`. C#: `Foo.cs` → `FooTests.cs` / `FooTest.cs` (incl. `*.Tests/` projects).
2. **Session is cache, source files are truth.** Every `tsma next` (and `tsma status`) re-scans the source and reconciles the function set: functions added or extracted since the last index surface as TODO, deleted functions are dropped, and existing progress is preserved. So a refactor that adds functions can never leave a stale "All functions complete!" — and if a test file is deleted, the function reverts to TODO regardless of what session.json says. Use `tsma rescan` to force this sync without touching progress (unlike `tsma reset --all`).
3. **Generated code is excluded.** `*_gen.go`, `*.pb.go` are not indexed.
4. **`.tsmignore` for custom exclusions.** Place a `.tsmignore` file in the project root to exclude paths from indexing. Same syntax as `.gitignore`.
5. **Coverage is measured, not enforced.** 100% is the goal, but unreachable branches (no DI, external dependencies) are accepted as DONE after one retry.

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

## .tsmignore

Place a `.tsmignore` file in the project root to exclude files and directories from indexing. Syntax matches `.gitignore`:

```
# Directories (trailing slash)
vendor/
internal/generated/

# File patterns
*.gen.go
*.pb.go

# Specific paths
cmd/legacy/main.go
```

| Pattern | Matches |
|---|---|
| `vendor/` | Directory named `vendor` at any depth |
| `*.gen.go` | Any file ending in `.gen.go` |
| `internal/db/` | The `internal/db` directory specifically |
| `cmd/main.go` | That exact file path |

If no `.tsmignore` exists, tsma uses built-in defaults only (`vendor/`, `.git/`, `.tsma/`, `node_modules/`).

## Language support

| Language | Detect marker | Indexer | Test runner | Coverage | Toolchain |
|---|---|---|---|---|---|
| Go | `go.mod` | `go/ast` | `go test` | `go test -coverprofile` | Go |
| TypeScript | `package.json` | regex | `npx vitest` / `npx jest` | `c8` / `istanbul` | Node.js |
| Python | `pyproject.toml` / `requirements.txt` / `setup.py` | regex | `pytest` | `coverage.py` | Python + pytest |
| Rust | `Cargo.toml` | regex | `cargo test` | `cargo llvm-cov` (llvm-cov) | cargo + `cargo-llvm-cov` (`llvm-tools-preview`) |
| Java | `pom.xml` / `build.gradle` / `build.gradle.kts` | regex | `mvn -Dtest=…` / `gradle test --tests`, module-aware (runs in the nearest submodule) | JaCoCo (`jacoco.xml`), read from the submodule's report | JDK + Maven or Gradle + JaCoCo plugin |
| C# | `*.csproj` / `*.sln` / `Directory.Build.props` | regex | `dotnet test --filter` | Cobertura via coverlet | .NET SDK + coverlet (`coverlet.collector`) |

## For LLM agents

Just run `tsma next`. The output tells the agent exactly what to do and when to run `tsma next` again. The loop continues until "All functions complete!".

## Reports

- [REPORT.md](REPORT.md) — tsma self-test (115 functions, PASS 87, DONE 28, coverage 95%)
- [REPORT.juicer.md](REPORT.juicer.md) — juicer project test (140 functions, PASS 114, DONE 26, coverage 96%)

## Caveat

100% branch coverage does not mean 100% correctness. tsma verifies that every branch is *exercised*, not that every assertion is *meaningful*. A test can hit all branches while checking nothing useful. Treat tsma results as a coverage floor, not a quality ceiling.

## Related work

[Diffblue Cover](https://www.diffblue.com/) uses reinforcement learning + symbolic execution to generate Java unit tests in a similar generate → verify → feedback loop. tsma takes the same core insight — deterministic verification driving iterative generation — but uses LLMs as the generator, works across Go/TypeScript/Python/Rust/Java/C#, and is open source.
