# Phase002 — submit 시 테스트 파일-함수 관계 검증

## 버그

`tsma submit Login ./contract_svc_test.go`처럼 **대상 함수와 무관한 테스트 파일**을 제출해도 DONE 처리된다.

## 근본 원인 (3건)

### 원인 1: 테스트 파일 위치/이름 검증 없음

`tsma submit`이 어떤 파일이든 받아들인다. 언어별 테스트 파일 관례를 검증하지 않는다.

| 언어 | 관례 | 예시 |
|---|---|---|
| Go | 같은 디렉토리, `_test.go` suffix | `handler.go` → `handler_test.go` (같은 디렉토리 필수) |
| TypeScript | 같은 디렉토리 또는 `__tests__/`, `.test.ts` 또는 `.spec.ts` | `handler.ts` → `handler.test.ts` |
| Python | 같은 디렉토리 또는 `tests/`, `test_` prefix | `auth_svc.py` → `test_auth_svc.py` |

### 원인 2: TotalBlocks == 0 → 100% 처리

`internal/coverage/compute_func_coverage.go:30`:

```go
if fc.TotalBlocks > 0 {
    fc.CoveredPct = float64(fc.CoveredBlocks) / float64(fc.TotalBlocks) * 100
} else {
    fc.CoveredPct = 100  // ← 커버리지 데이터가 없으면 100%로 간주
}
```

잘못된 패키지의 테스트 파일을 제출하면, 커버리지 프로파일에 대상 함수의 블록이 아예 없다. `TotalBlocks == 0`이 되고, "빈 함수"로 간주하여 100%를 반환한다.

### 원인 3: 패키지 전체 테스트 실행 (Go)

`internal/coverage/go_checker_check.go`에서 `-run` 플래그 없이 패키지 전체 테스트를 실행한다. 같은 패키지 내 다른 테스트 파일이 대상 함수를 간접적으로 호출하면, 제출한 테스트와 무관하게 커버리지가 잡힌다.

## 수정

### 수정 1: 테스트 파일 위치/이름 검증 (submit 진입 시)

`internal/cli/submit.go`에서 커버리지 실행 전에 검증. 실패 시 즉시 거부.

**Go**: 테스트 파일이 대상 함수와 같은 디렉토리에 있어야 한다.
```go
testDir := filepath.Dir(testFile)
funcDir := filepath.Dir(fn.File)
if testDir != funcDir {
    return fmt.Errorf("test file must be in the same directory as the function\n  test: %s\n  func: %s", testDir, funcDir)
}
```

**TypeScript**: 테스트 파일명이 `.test.ts`, `.test.js`, `.spec.ts`, `.spec.js`로 끝나야 한다.
```go
base := filepath.Base(testFile)
if !strings.Contains(base, ".test.") && !strings.Contains(base, ".spec.") {
    return fmt.Errorf("test file must have .test.ts or .spec.ts suffix: %s", base)
}
```

**Python**: 테스트 파일명이 `test_` prefix로 시작해야 한다.
```go
base := filepath.Base(testFile)
if !strings.HasPrefix(base, "test_") {
    return fmt.Errorf("test file must have test_ prefix: %s", base)
}
```

### 수정 2: TotalBlocks == 0 → 0% 처리

`internal/coverage/compute_func_coverage.go`:

```go
if fc.TotalBlocks > 0 {
    fc.CoveredPct = float64(fc.CoveredBlocks) / float64(fc.TotalBlocks) * 100
} else {
    fc.CoveredPct = 0  // 커버리지 데이터 없음 = 미커버
}
```

TS/Python checker의 동일 패턴도 수정:
- `internal/coverage/compute_ts_func_coverage.go`
- `internal/coverage/compute_py_func_coverage.go`

### 수정 3: 제출된 테스트 파일의 테스트 함수만 실행 (Go)

`internal/coverage/go_checker_check.go`에서 `-run` 플래그 추가:

```go
testFuncs := runner.ExtractTestFuncs(testFile)
args := []string{"test", "-count=1",
    "-coverprofile=" + coverFile,
    "-covermode=set",
}
if len(testFuncs) > 0 {
    args = append(args, "-run", strings.Join(testFuncs, "|"))
}
args = append(args, pkgPath)
```

TS/Python은 이미 파일 단위 실행(`npx vitest run <file>`, `pytest <file>`)이므로 수정 불필요. **Go만 수정**.

`extractTestFuncs`는 현재 `internal/runner/` 패키지의 비공개 함수이므로, `ExtractTestFuncs`로 공개 전환.

## 변경 파일

| 파일 | 변경 |
|---|---|
| `internal/cli/submit.go` | 수정 1: 테스트 파일 위치/이름 검증 추가 (언어별 분기) |
| `internal/coverage/compute_func_coverage.go:30` | 수정 2: `CoveredPct = 100` → `CoveredPct = 0` |
| `internal/coverage/compute_ts_func_coverage.go` | 수정 2: 동일 패턴 |
| `internal/coverage/compute_py_func_coverage.go` | 수정 2: 동일 패턴 |
| `internal/coverage/go_checker_check.go` | 수정 3: `-run` 플래그 추가 |
| `internal/runner/extract_test_funcs.go` | 수정 3: `extractTestFuncs` → `ExtractTestFuncs` 공개 전환 |

## 검증

```bash
# 수정 1 검증: 다른 디렉토리 파일 → 즉시 거부
tsma submit "internal/api/contract.Handler.ListContracts" ./internal/service/contract_svc_test.go
# 기대: Error: test file must be in the same directory as the function

# 수정 2 검증: 같은 디렉토리지만 대상 함수의 블록이 프로파일에 없는 경우 → 0%
# (엣지 케이스: 빈 함수나 프로파일 누락)

# 수정 3 검증: 같은 디렉토리, 올바른 파일
tsma reset "internal/api/contract.Handler.ListContracts"
tsma submit "internal/api/contract.Handler.ListContracts" ./internal/api/contract/handler_test.go
# 기대: 제출 파일의 Test* 함수만 실행한 커버리지 (65%)

# 전체 빌드 + 테스트
go build ./...
go test ./... -count=1
filefunc validate .
```

## 의존성

없음. cli + coverage + runner 내부 변경.
