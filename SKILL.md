---
name: tsma
description: Test coverage ratchet for Go, TypeScript, and Python. Indexes every function, detects missing tests, measures branch coverage, and guides LLM agents to fill gaps one function at a time. Use this skill when writing unit tests for legacy codebases, measuring test coverage, or running test-generation loops with LLM agents. Triggers on tasks involving tsma commands, test coverage improvement, or missing test detection.
license: MIT
metadata:
  author: park-jun-woo
  version: "1.0.0"
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

Requires Go 1.25+.

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
  ▶ Cover the uncovered branches. Next tsma next will re-measure.
```

The agent sees exactly which lines to cover — no guessing.

## Why Some Functions Can't Reach 100%

Functions with **interface dependencies** (mockable) can reach 100%. Functions with **concrete struct dependencies** (not mockable) are stuck below 100% because you can't control internal dependency behavior in unit tests. tsma marks these as DONE after one retry.

## Language Support

| Language | Indexer | Test runner | Coverage |
|---|---|---|---|
| Go | `go/ast` | `go test` | `go test -coverprofile` |
| TypeScript | regex | `npx vitest` / `npx jest` | `c8` / `istanbul` |
| Python | regex | `pytest` | `coverage.py` |

## Agent Instructions

Give this to your LLM agent:

```
1. Run `tsma next`
2. If TODO — read the function, write a test
3. If test fails — read the error, fix the test
4. If uncovered branches shown — add tests for those branches
5. If PASS/DONE — next function is automatically shown
6. Repeat until "All functions complete!"
```

## Validated Results

| Project | Functions | PASS (100%) | DONE | Coverage |
|---|---|---|---|---|
| tsma self-test | 115 | 87 | 28 | 95% |
| juicer | 140 | 114 | 26 | 96% |
| gozhip (527 funcs) | 527 | 246 | 281 | — |
