# tsma

[![Version](https://img.shields.io/badge/version-v0.5.1-blue.svg)](https://github.com/park-jun-woo/tsma/releases)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![skills.sh](https://skills.sh/b/park-jun-woo/tsma)](https://skills.sh/park-jun-woo/tsma)

Regression defense for legacy code. Built on the [reins](https://github.com/park-jun-woo/reins) deterministic quest gate. Indexes every function, detects missing tests, measures branch coverage, and drives LLM agents to fill the gaps — one function at a time or in an unattended loop.

## What tsma does

1. **Shows where tests are missing.** Indexes every function, matches against test files by convention, and tells you exactly which functions have no tests.
2. **Uncovered branch feedback makes LLM tests dramatically better.** Without feedback, LLM writes 60-70% coverage. With specific line numbers ("line 41, 44, 70 uncovered"), LLM reaches 100% in one shot.
3. **Unattended loop for LLM agents.** `tsma scan .` indexes functions, then `tsma loop` drives the full generate / gate / retry cycle — no human interaction needed. For manual control, use `tsma next` and `tsma submit`.

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

### Unattended (recommended)

```bash
tsma scan .           # index all functions, seed TODO items
tsma loop             # LLM generates tests, gate verifies, retries on FAIL
```

The loop picks each TODO function, generates a test via LLM (`claude:sonnet` by default), runs the test, measures branch coverage, and locks PASS at 100%. On FAIL it feeds back uncovered lines and retries (up to 3 attempts per function).

### Manual (agent-driven)

```bash
tsma scan .                            # index all functions (once)
tsma next                              # shows next TODO + authoring prompt
# write the test file to disk
tsma submit --key <FunctionName> --in <test_file>
# PASS (100%) → locked, or FAIL (with uncovered line feedback)
tsma next                              # next function
```

Repeat `next` / write / `submit` until `tsma status` shows all complete.

## Commands

| Command | Purpose |
|---|---|
| `tsma scan [dir]` | Index functions in the project, seed TODO items (default `.`) |
| `tsma next` | Show the next TODO function with authoring prompt |
| `tsma submit --key <key> --in <file>` | Submit a test; gate evaluates and returns a verdict |
| `tsma status` | Progress tally (TODO/PASS/REVIEW/DONE) |
| `tsma export` | Emit terminal results as JSONL |
| `tsma rules` | Show the gate's violation-rule catalog |
| `tsma loop [--model backend:model] [--max-items N]` | Unattended generate / gate / retry loop |

## Status

| Status | Meaning |
|---|---|
| **TODO** | Not yet verified — needs a test |
| **PASS** | 100% branch coverage, locked by gate (irreversible) |
| **REVIEW** | 100% coverage but test uses escape hatches (unsafe, reflect, etc.) — needs human review |
| **DONE** | Max retries (3) reached without PASS — auto-accepted |

## Example session

Scan the project to seed TODO items:

```
$ tsma scan .
Detected: go
Indexed 527 functions
```

View the next TODO function:

```
$ tsma next
Login  TODO
  Source: internal/api/auth/handler.go:86-119
```

After writing a test and submitting:

```
$ tsma submit --key Login --in internal/api/auth/handler_test.go
PASS  Login  (100% branch coverage)
```

When coverage is not 100%, the verdict carries uncovered line locations:

```
$ tsma submit --key ListContracts --in internal/api/contract/handler_test.go
FAIL  branch-coverage-below-100
  where: handler.go:41, handler.go:44, handler.go:70
  expected: 100% branch coverage
  actual: 65.0% (3 uncovered branch(es))
```

The agent sees exactly which lines to cover, fixes the test, and submits again.

## Principles

1. **Convention-based test matching.** Go: `handler.go` → `handler_test.go` in the same directory. TS: `.test.ts` / `.spec.ts`. Python: `test_` prefix. Rust: in-file `#[cfg(test)] mod tests` or `tests/*.rs`. Java: `src/main/java/…/Foo.java` → `src/test/java/…/FooTest.java`. C#: `Foo.cs` → `FooTests.cs` / `FooTest.cs` (incl. `*.Tests/` projects).
2. **Deterministic gate, irreversible ratchet.** Only the machine locks PASS (authority asymmetry). Once PASS, it is irreversible — remaining work monotonically decreases. The gate evaluates from disk truth: matched tests are re-run and re-measured on every submit.
3. **Generated code is excluded.** `*_gen.go`, `*.pb.go` are not indexed.
4. **`.tsmignore` for custom exclusions.** Place a `.tsmignore` file in the project root to exclude paths from indexing. Same syntax as `.gitignore`.
5. **Coverage is measured, not enforced.** 100% is the goal, but unreachable branches (no DI, external dependencies) are accepted as DONE after 3 retries.
6. **Escape-hatch detection.** Tests that reach 100% coverage by cheating (unsafe, reflect, linkname, `as any`, etc.) are flagged REVIEW for human inspection rather than silently locked as PASS.

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

**tsma's response:** After 3 retries with uncovered branch feedback, tsma marks the function as DONE with the achieved coverage. This is not a tool limitation — it reflects the code's testability. Introducing an interface (DI) would make the function fully testable, but that requires modifying the source code.

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

**Unattended:** `tsma scan .` then `tsma loop` — the loop handles generation, verification, retry, and progress automatically.

**Manual:** `tsma scan .` then repeat `tsma next` / write test / `tsma submit --key <key> --in <file>` until `tsma status` shows all complete.

## Reports (pre-reins runs)

- [REPORT.md](REPORT.md) — tsma self-test (115 functions, PASS 87, DONE 28, coverage 95%)
- [REPORT.juicer.md](REPORT.juicer.md) — juicer project test (140 functions, PASS 114, DONE 26, coverage 96%)
- [REPORT.gozhip.md](REPORT.gozhip.md) — gozhip project (527 functions)

## Caveat

100% branch coverage does not mean 100% correctness. tsma verifies that every branch is *exercised*, not that every assertion is *meaningful*. A test can hit all branches while checking nothing useful. Treat tsma results as a coverage floor, not a quality ceiling.

## Related work

[Diffblue Cover](https://www.diffblue.com/) uses reinforcement learning + symbolic execution to generate Java unit tests in a similar generate → verify → feedback loop. tsma takes the same core insight — deterministic verification driving iterative generation — but uses LLMs as the generator, works across Go/TypeScript/Python/Rust/Java/C#, and is open source.
