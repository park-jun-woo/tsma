---
name: tsma
description: Test coverage ratchet for Go, TypeScript, Python, Rust, Java, and C#. Built on the reins deterministic quest gate, indexes every function, detects missing tests, measures branch coverage, and drives LLM agents to fill gaps — one function at a time or in an unattended loop. Use this skill when writing unit tests for legacy codebases, measuring test coverage, or running test-generation loops with LLM agents. Triggers on tasks involving tsma commands, test coverage improvement, or missing test detection.
license: MIT
metadata:
  author: park-jun-woo
  version: "0.5.1"
---

# tsma — Test Coverage Ratchet for LLM Agents

tsma indexes every function in a codebase, detects missing tests, measures branch coverage, and drives an LLM agent to fill the gaps. Built on the [reins](https://github.com/park-jun-woo/reins) deterministic quest gate — generation is probabilistic, verification is deterministic. Once a function reaches 100% branch coverage, PASS is locked and irreversible.

## When to Use This Skill

- Improving test coverage on a legacy codebase
- Writing unit tests guided by branch coverage feedback
- Running an unattended test-generation loop (`tsma loop`)
- Measuring which functions lack tests

## Install

```bash
go install github.com/park-jun-woo/tsma/cmd/tsma@latest
```

**Prerequisites:** Go 1.22+. Optional: `tree-sitter` CLI for enhanced indexing (set `TSMA_TREE_SITTER` to override the binary path; without it tsma falls back to line-based indexing).

## Quick Start

### Unattended (recommended)

```bash
tsma scan .           # index all functions, seed TODO items
tsma loop             # LLM generates tests, gate verifies, retries on FAIL
```

The loop picks each TODO function, generates a test via LLM (`claude:sonnet` by default), runs the test, measures branch coverage, and locks PASS at 100%. On FAIL it feeds back uncovered lines and retries (up to 3 attempts per function). No human interaction needed.

### Manual (agent-driven)

```bash
tsma scan .                            # index all functions (once)
tsma next                              # shows next TODO + authoring prompt
# agent writes the test file to disk
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
| `tsma submit --key <key> --in <file>` | Submit a test; gate evaluates and returns a verdict (PASS/FAIL/REVIEW) |
| `tsma status` | Progress tally (TODO/PASS/REVIEW/DONE) |
| `tsma export` | Emit terminal results as JSONL (emit-once) |
| `tsma rules` | Show the gate's violation-rule catalog |
| `tsma loop [--model backend:model] [--max-items N]` | Unattended generate / gate / retry loop |

## Status Types

| Status | Meaning |
|---|---|
| **TODO** | Not yet verified — needs a test |
| **PASS** | 100% branch coverage, locked by gate (irreversible) |
| **REVIEW** | 100% coverage but test uses escape hatches (unsafe, reflect, etc.) — needs human review |
| **DONE** | Max retries (3) reached without PASS — auto-accepted |

## Gate Rules

tsma's gate has two severity levels. Run `tsma rules` to see the full catalog.

**Fail rules** (block PASS):

| Rule ID | Fires when |
|---|---|
| `tests-must-pass` | Tests don't compile, fail, or coverage measurement errors |
| `branch-coverage-below-100` | Tests pass but branch coverage < 100% |

**Review rules** (surface only when all Fail rules are clean — tests pass at 100%):

| Rule ID | Language | Detects |
|---|---|---|
| `TS-REFL-001` | Go | `unsafe` package usage in tests |
| `TS-REFL-002` | Go | Dynamic `reflect` (MethodByName/FieldByName) in tests |
| `TS-REFL-003` | Go | `//go:linkname` in tests |
| `TS-REFL-TS-001` | TypeScript | `as any` casts in tests |
| `TS-REFL-TS-002` | TypeScript | Reflect API usage in tests |
| `TS-REFL-TS-003` | TypeScript | `Object.getOwnProperty*` in tests |
| `TS-REFL-JV-001` | Java | `java.lang.reflect getDeclared*` in tests |
| `TS-REFL-JV-002` | Java | `setAccessible(true)` in tests |
| `TS-REFL-CS-001` | C# | `System.Reflection` GetMethod/Field/Property in tests |
| `TS-REFL-CS-002` | C# | MethodInfo/PropertyInfo/FieldInfo declarations in tests |
| `TS-REFL-RS-001` | Rust | `unsafe` blocks in tests |
| `TS-REFL-RS-002` | Rust | `transmute` in tests |
| `TS-REFL-RS-003` | Rust | Raw pointer (`std::ptr`) access in tests |

## Key Insight: Coverage Feedback Changes Everything

Without feedback, LLMs write 60-70% coverage tests. With specific uncovered line numbers, LLMs reach 100% in one shot.

When `submit` returns FAIL with the `branch-coverage-below-100` rule, the verdict carries a Fact with exactly which lines are uncovered:

```
FAIL  branch-coverage-below-100
  where: handler.go:41, handler.go:44, handler.go:70
  expected: 100% branch coverage
  actual: 65.0% (3 uncovered branch(es))
```

The agent sees exactly which lines to cover — no guessing.

## Loop Command (Unattended Drive)

`tsma loop` closes the generate / gate / retry cycle in-process:

```bash
tsma scan .
tsma loop                              # uses claude:sonnet (CLI login, no API key)
tsma loop --model ollama:gemma4        # use a local model
tsma loop --max-items 10               # limit to 10 functions
```

The loop for each TODO function:
1. Renders the authoring prompt (function body, existing tests, package context)
2. LLM generates a test file
3. Gate evaluates: run tests, measure coverage, apply rules
4. PASS: lock and advance to next function
5. FAIL: feed back uncovered lines + rule-specific coaching, retry (up to 3 attempts)
6. 3 FAILs: lock DONE, move on

**Authority asymmetry:** the LLM only generates. Only the deterministic gate locks PASS. The loop calls the same gate path as `submit` — feedback parity is guaranteed.

**Supported backends:** `claude:<model>`, `ollama:<model>`, `xai:<model>`, `gemini:<model>`, `grok:<model>`, `codex:<model>`, `geminicli:<model>`. See the [reins manual](https://github.com/park-jun-woo/reins) for backend details.

## Why Some Functions Can't Reach 100%

Functions with **interface dependencies** (mockable) can reach 100%. Functions with **concrete struct dependencies** (not mockable) are stuck below 100% because you can't control internal dependency behavior in unit tests. tsma marks these as DONE after 3 retries.

## Language Support

| Language | Detect marker | Test runner | Coverage | Toolchain |
|---|---|---|---|---|
| Go | `go.mod` | `go test` | `go test -coverprofile` | Go |
| TypeScript | `package.json` | `npx vitest` / `npx jest` | `c8` / `istanbul` | Node.js |
| Python | `pyproject.toml` / `requirements.txt` / `setup.py` | `pytest` | `coverage.py` | Python + pytest |
| Rust | `Cargo.toml` | `cargo test` | `cargo llvm-cov` | cargo + `cargo-llvm-cov` (`llvm-tools-preview`) |
| Java | `pom.xml` / `build.gradle(.kts)` | `mvn` / `gradle test` | JaCoCo | JDK + Maven or Gradle + JaCoCo plugin |
| C# | `*.csproj` / `*.sln` / `Directory.Build.props` | `dotnet test` | Cobertura via coverlet | .NET SDK + coverlet |

Indexers are AST-based for Go and regex-based for all other languages. Tree-sitter enhances indexing when available.

## IMPORTANT: Do NOT modify .tsmignore

**NEVER add files or directories to `.tsmignore` to skip writing tests.** `.tsmignore` exists solely for the project owner to exclude generated code or vendored dependencies. Using it to avoid test writing defeats the entire purpose of tsma. If a function is hard to test, write the best test you can — tsma will mark it DONE after 3 retries. Do not circumvent the process.

## Agent Instructions

### Unattended mode

```bash
tsma scan .
tsma loop
```

The loop handles everything automatically — generation, verification, retry, and progress.

### Manual mode

1. `tsma scan .` — seed the session (once)
2. `tsma next` — read the authoring prompt for the next TODO function
3. Write the test file to disk
4. `tsma submit --key <key> --in <test_file>` — gate evaluates from disk
5. If FAIL: read the Fact feedback (uncovered lines, compile errors), fix the test, submit again
6. If PASS: run `tsma next` for the next function
7. Repeat until `tsma status` shows all functions complete
