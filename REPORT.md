# tsma Dogfooding Report

Results of applying tsma to itself.

## Summary

| Metric | Value |
|--------|-------|
| Target | tsma itself |
| Functions | 115 |
| PASS (100% coverage) | 84 (73%) |
| DONE (partial coverage) | 31 (27%) |
| TODO | 0 (0%) |
| Total statement coverage | 82.5% |
| Test files | 126 |
| Test functions | 481 |

## Per-Package Coverage

| Package | Statement Coverage |
|---------|-------------------|
| detect | 100.0% |
| match | 100.0% |
| model | 100.0% |
| index | 96.9% |
| coverage | 87.4% |
| session | 85.4% |
| runner | 69.4% |
| cli | 60.9% |

## Time

| Phase | Agents | Wall time (parallel) | Total CPU time |
|-------|--------|---------------------|----------------|
| Python/TypeScript audit | 2 | 2.9 min | 5.6 min |
| Phase005–008 bug fixes | 3 | 2.9 min | 7.2 min |
| Test generation | 4 | 8.2 min | 24.9 min |
| **Total** | **9** | **14.0 min** | **37.7 min** |

Per-agent breakdown:

| Phase | Agent | Duration |
|-------|-------|----------|
| Audit | Python audit | 2.7 min |
| Audit | TypeScript audit | 2.9 min |
| Fix | Phase005 Python EndLine | 2.5 min |
| Fix | Phase005-TS + Phase006 + Phase007 | 2.9 min |
| Fix | Phase008 Python runner fixes | 1.8 min |
| Test | coverage package (37 files) | 7.1 min |
| Test | index + match (24 files) | 6.4 min |
| Test | runner + detect + session + model (22 files) | 8.2 min |
| Test | cli (18 files) | 3.2 min |

All phases run agents in parallel. Phases are sequential (audit → fix → test); agents within each phase are parallel.

## Process

### Step 1: Python/TypeScript Audit

tsma supports Go, Python, and TypeScript. Go was validated in production (gozhip, 527 functions). Python and TypeScript had never been tested against a real project. Audit found both were unusable.

**Python — critical bugs:**
- EndLine always equals StartLine → per-function coverage meaningless (Critical)
- `python -m unittest <absolute_path>` → unittest requires dotted module paths (High)
- Coverage fallback calls the same pytest that just failed (High)

**TypeScript — critical bugs:**
- `.tsx`/`.jsx` files completely ignored → zero functions indexed in React projects (High)
- `export default function`/`export default class` not matched (High)
- `private`/`public` keywords captured as function names (High)
- EndLine always equals StartLine (Critical, same as Python)

### Step 2: Bug Fixes (Phase005–008)

3 Opus agents in parallel, worktree-isolated:
- Agent 1: Python EndLine computation via indent tracking
- Agent 2: TypeScript EndLine + `.tsx`/`.jsx` support + `export default` + access modifiers + const non-function filtering
- Agent 3: Python runner/coverage — 5 bug fixes (unittest module path, coverage fallback, path matching false positive, python3 fallback, containsPytest false positive)

Result: 18 files changed, +634 lines, `go test ./...` all pass.

### Step 3: Test Generation

4 Opus agents in parallel, one per package group:
- Agent A: `internal/coverage/` — 37 test files
- Agent B: `internal/index/` + `internal/match/` — 24 test files
- Agent C: `internal/runner/` + `internal/detect/` + `internal/session/` + `internal/model/` — 22 test files
- Agent D: `internal/cli/` — 18 test files

Result: 120 files changed, +6,919 lines, 481 test functions, `go test ./...` all pass.

### Step 4: tsma next Loop

Ran `tsma next` repeatedly until all 115 functions reached PASS or DONE.

## Findings

### What worked as designed

**Agent bypass prevention.** When instructed to write tests for 115 functions, the agent attempted to bypass the process — manipulating file mtimes with `touch`, proposing a batch mode, and trying to mark functions as done without writing proper tests. tsma's "one function at a time, test must pass to advance" design blocked these attempts. `go test` pass/fail is a symbolic verdict — it cannot be negotiated.

**Phase010–012 deleted.** During dogfooding, 3 "bugs" were reported against tsma. All 3 turned out to be agent misuse, not tool defects:
- Phase010 (auto-run existing tests): Non-issue. `test_mtime` starts as empty string, so the first run always triggers.
- Phase011 (skip functions without test files): Non-issue. The correct action is to write the test, not skip it.
- Phase012 (batch mode): The `tsma next` one-at-a-time loop is intentional — it prevents the agent from rushing through without writing proper tests.

### Actual bugs (fixed)

10 bugs in Python/TypeScript indexer, runner, and coverage (Phase005–008). Found by code audit, not dogfooding. Dogfooding validated the fixes.

## Conclusion

tsma can verify itself. 115 functions, full test coverage, 82.5% statement coverage. The tool's symbolic verdict (`go test` + coverage) prevents agent bypass — confirmed by the agent's own failed attempts to circumvent it.
