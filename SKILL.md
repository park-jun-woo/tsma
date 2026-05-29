---
name: tsma
description: Test coverage ratchet for Go, TypeScript, Python, Rust, Java, and C#. Indexes every function, detects missing tests, measures branch coverage, and guides LLM agents to fill gaps one function at a time. Use this skill when writing unit tests for legacy codebases, measuring test coverage, or running test-generation loops with LLM agents. Triggers on tasks involving tsma commands, test coverage improvement, or missing test detection.
license: MIT
metadata:
  author: park-jun-woo
  version: "0.3.0"
---

# tsma — Test Coverage Ratchet for LLM Agents

tsma indexes every function in a codebase, detects missing tests, measures branch coverage, and drives an LLM agent to fill the gaps — one function at a time.

## When to Use This Skill

- Improving test coverage on a legacy codebase
- Writing unit tests guided by branch coverage feedback
- Running a test-generation loop (`tsma next` until all functions covered)
- Measuring which functions lack tests

## Install

```bash
go install github.com/park-jun-woo/tsma/cmd/tsma@latest
```

**Prerequisites:** Go 1.25+ and gcc (cgo required). tsma uses cgo dependencies — without gcc the build fails.

## The Only Command You Need

```bash
tsma next
```

This single command drives the entire workflow:

```
tsma next  → shows next function without a test (TODO)
           → agent writes a test
tsma next  → runs test, measures coverage
           → 100%? PASS, moves to next function
           → <100%? shows uncovered branch line numbers
           → agent adds tests for uncovered branches
tsma next  → re-measures, marks DONE or PASS, moves on
```

Repeat until `All functions complete!`.

## Other Commands

| Command | Purpose |
|---|---|
| `tsma next` | Full workflow: detect → test → coverage → advance |
| `tsma list [--page N]` | All functions with status |
| `tsma status` | Progress summary |
| `tsma reset --all` | Delete session |

## Status Types

| Status | Meaning |
|---|---|
| **TODO** | No test file, or coverage < 100% |
| **DONE** | Not 100% but best effort (unreachable branches) |
| **PASS** | 100% branch coverage |

## Key Insight: Coverage Feedback Changes Everything

Without feedback, LLM writes 60-70% coverage tests. With specific uncovered line numbers, LLM reaches 100% in one shot.

```
tsma next

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

The agent sees exactly which lines to cover — no guessing.

## Why Some Functions Can't Reach 100%

Functions with **interface dependencies** (mockable) can reach 100%. Functions with **concrete struct dependencies** (not mockable) are stuck below 100% because you can't control internal dependency behavior in unit tests. tsma marks these as DONE after one retry.

## Language Support

| Language | Detect marker | Test runner | Coverage | Toolchain |
|---|---|---|---|---|
| Go | `go.mod` | `go test` | `go test -coverprofile` | Go |
| TypeScript | `package.json` | `npx vitest` / `npx jest` | `c8` / `istanbul` | Node.js |
| Python | `pyproject.toml` / `requirements.txt` / `setup.py` | `pytest` | `coverage.py` | Python + pytest |
| Rust | `Cargo.toml` | `cargo test` | `cargo llvm-cov` | cargo + `cargo-llvm-cov` (`llvm-tools-preview`) |
| Java | `pom.xml` / `build.gradle(.kts)` | `mvn` / `gradle test` | JaCoCo | JDK + Maven or Gradle + JaCoCo plugin |
| C# | `*.csproj` / `*.sln` / `Directory.Build.props` | `dotnet test` | Cobertura via coverlet | .NET SDK + coverlet |

Indexers are AST-based for Go and regex-based for all other languages.

## IMPORTANT: Do NOT modify .tsmignore

**NEVER add files or directories to `.tsmignore` to skip writing tests.** `.tsmignore` exists solely for the project owner to exclude generated code or vendored dependencies. Using it to avoid test writing defeats the entire purpose of tsma. If a function is hard to test, write the best test you can — tsma will mark it DONE after one retry. Do not circumvent the process.

## Agent Instructions

Just run `tsma next`. The output tells the agent exactly what to do and when to run `tsma next` again. The loop continues until "All functions complete!".

## Validated Results

| Project | Functions | PASS (100%) | DONE | Coverage |
|---|---|---|---|---|
| tsma self-test | 115 | 87 | 28 | 95% |
| juicer | 140 | 114 | 26 | 96% |
| gozhip (527 funcs) | 527 | 246 | 281 | — |
