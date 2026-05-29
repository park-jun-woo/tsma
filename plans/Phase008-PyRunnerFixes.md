# Phase008 — Python 러너·커버리지 버그 수정

## 문제

### 1. unittest 실행 시 파일 경로 전달 — 실행 실패

`py_runner_run.go:23`에서 `python -m unittest <절대경로>` 실행. unittest는 파일 경로를 받지 않고 점 구분 모듈 경로(`tests.test_handler`)를 요구. `ModuleNotFoundError` 발생.

### 2. 커버리지 폴백이 또 pytest 사용

`run_coverage_py.go:12`에서 폴백 명령이 `python -m coverage run -m pytest <file>`. pytest 실패 시 호출되는 폴백인데 동일하게 pytest 사용 → 이중 실패.

### 3. 커버리지 경로 매칭 false positive

`matches_py_path.go:23`에서 `strings.HasSuffix(normalized, normalizedTarget)`. `other_handler.py`가 `handler.py`에 매칭되어 다른 파일의 커버리지를 가져옴.

### 4. `filepath.Abs` CWD 기준 해석

`py_runner_run.go:12`에서 `filepath.Abs(testFile)`가 현재 작업 디렉토리 기준. `projectRoot`와 다른 위치에서 tsma 실행 시 경로 오류.

### 5. `python` 바이너리 하드코딩

많은 Linux 배포판에서 `python`이 없고 `python3`만 존재.

### 6. `containsPytest` false positive

`contains_pytest.go:17`에서 `pattern` 파라미터와 별개로 항상 `strings.Contains(content, "pytest")`를 체크. `# not using pytest` 같은 주석만 있어도 true 반환.

## 수정 방향

### unittest 모듈 경로 변환

파일 경로를 모듈 경로로 변환:

```go
// /home/user/project/tests/test_handler.py
// projectRoot = /home/user/project
// → tests.test_handler

rel, _ := filepath.Rel(projectRoot, absTest)
modulePath := strings.TrimSuffix(rel, ".py")
modulePath = strings.ReplaceAll(modulePath, string(filepath.Separator), ".")
// cmd: python -m unittest modulePath -v
```

### 커버리지 폴백

pytest 폴백 대신 unittest 사용:

```go
// 현재: python -m coverage run -m pytest <file>
// 변경: python -m coverage run -m unittest <modulePath> -v
```

### 경로 매칭 수정

suffix 매칭 전에 경로 구분자 확인:

```go
func matchesPyPath(coveragePath, targetFile string) bool {
    // ... 기존 정규화 ...
    if normalized == normalizedTarget {
        return true
    }
    // suffix 매칭 시 앞 문자가 / 인지 확인
    if strings.HasSuffix(normalized, normalizedTarget) {
        idx := len(normalized) - len(normalizedTarget)
        if idx > 0 && normalized[idx-1] == '/' {
            return true
        }
    }
    return false
}
```

### filepath.Abs 수정

`filepath.Abs` 대신 `filepath.Join(projectRoot, testFile)` 사용:

```go
absTest := filepath.Join(projectRoot, testFile)
```

### python3 폴백

`python3` 먼저 시도, 없으면 `python`:

```go
func findPython() string {
    if _, err := exec.LookPath("python3"); err == nil {
        return "python3"
    }
    return "python"
}
```

### containsPytest 수정

`pattern` 파라미터의 정확한 매칭만 사용. 무조건 `"pytest"` 서브스트링 체크하는 fallback 제거:

```go
// 현재
return strings.Contains(content, strings.ToLower(pattern)) || strings.Contains(content, "pytest")
// 변경
return strings.Contains(content, strings.ToLower(pattern))
```

## 변경 파일

| 파일 | 변경 |
|---|---|
| `internal/runner/py_runner_run.go` | unittest 모듈 경로 변환, `filepath.Join` 사용, `python3` 폴백 |
| `internal/coverage/run_coverage_py.go` | 폴백을 unittest 기반으로 변경 |
| `internal/coverage/matches_py_path.go` | suffix 매칭에 경로 구분자 확인 추가 |
| `internal/coverage/run_pytest_cov.go` | `filepath.Join` 사용 |
| `internal/runner/contains_pytest.go` | 무조건 `"pytest"` 서브스트링 체크 제거 |

## 검증

```bash
go test ./internal/runner/... -count=1 -run "Py"
go test ./internal/coverage/... -count=1 -run "Py"

# 추가 테스트 케이스
# - pytest 없는 환경에서 unittest 폴백
# - python3만 있는 환경에서 바이너리 탐색
# - __init__.py 없는 디렉토리의 모듈 경로 변환
# - containsPytest: "# not using pytest" 주석만 있는 파일 → false
```

## 의존성

Phase005(EndLine) 선행 권장. EndLine 없으면 커버리지 수정해도 1줄만 측정.
