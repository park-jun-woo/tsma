# Phase004 — 파일 기반 테스트 매칭으로 전환

## 문제

함수명 기반 매칭(`TestLogin` → `Login`)은 추측이다. `TestLogin`이 `Handler.Login`을 테스트하는지 `strictHandler.Login`을 테스트하는지 이름만으로 알 수 없다. gozhip admin에서 생성 코드(server.gen.go)의 74개 함수가 오매칭으로 DONE 처리되었다.

## 정적으로 확정 가능한 사실

1. `_test.go` 파일이 존재한다
2. 그 안에 `func Test*`가 있다
3. `go test`가 pass/fail한다

어떤 함수가 실제로 테스트되는지는 커버리지(`tsma cover`)를 돌려야만 알 수 있다.

## 수정 방향

매칭 단위를 함수 → **파일**로 변경.

### 상태 판정 (변경)

| 조건 | 상태 |
|---|---|
| 소스 파일에 대응하는 테스트 파일 없음 | **TODO** |
| 테스트 파일 있으나 `go test` 실패 | **FAIL** |
| 테스트 파일 있고 `go test` 성공 | **DONE** |

DONE = "대응하는 테스트 파일이 존재하고 패키지 테스트가 pass한다". 해당 파일의 모든 함수가 실제로 커버되는지는 보장하지 않는다 — 파일 단위 매칭의 한계. 함수별 커버리지 확인은 `tsma cover`의 영역.

### 파일 매칭 규칙

| 언어 | 소스 | 테스트 | 규칙 |
|---|---|---|---|
| Go | `handler.go` | `handler_test.go` | 같은 디렉토리, `_test.go` suffix |
| TS | `handler.ts` | `handler.test.ts` / `handler.spec.ts` | 같은 디렉토리 또는 `__tests__/` |
| Python | `auth_svc.py` | `test_auth_svc.py` | 같은 디렉토리 또는 `tests/` |

함수명 매칭 없음. 파일명만으로 판정.

### 생성 코드 제외

`*_gen.go`, `*.gen.go`, `*.pb.go` — 인덱싱에서 제외. 코드젠 도구가 보장하는 영역이며 수동 테스트 대상이 아니다.

## 변경 사항

### 1. `internal/match/` — 파일 기반으로 단순화

현재: `Match(projectRoot, *Function) → (testFile, found)` — 함수명으로 테스트 함수 탐색.

변경: `Match(projectRoot, sourceFile string) → (testFile string, found bool)` — 소스 파일에 대응하는 테스트 파일 존재 여부만 확인.

삭제할 파일:
- `matches_test_func.go` — 함수명 매칭
- `contains_test_for.go` — Test* 함수 파싱
- `ts_test_mentions_func.go` — describe/test 패턴 매칭
- `py_test_mentions_func.go` — def test_ 패턴 매칭

### 2. `internal/index/` — 생성 코드 제외

Go: `*_gen.go`, `*.gen.go`, `*.pb.go` 스킵 추가.

### 3. `internal/cli/check.go` — 파일 기반 판정

현재: 함수별로 match 호출.
변경: 소스 파일별로 match 호출. 테스트 파일이 있는 소스 파일의 모든 함수 → TESTED (pass 시) / FAIL (fail 시).

## 변경 파일

| 파일 | 변경 |
|---|---|
| `internal/match/` 전체 | 파일 기반으로 재작성, 함수명 매칭 관련 파일 삭제 |
| `internal/index/go_indexer_index.go` | `*_gen.go`, `*.gen.go`, `*.pb.go` 스킵 |
| `internal/index/is_go_source.go` | 생성 파일 제외 조건 추가 |
| `internal/cli/check.go` | 파일 기반 매칭 + 패키지 단위 테스트 |
| `internal/cli/analyze_project.go` | 파일 기반 매칭 호출 |

model, status, summary는 변경 없음 (StatusDone 유지).

## 검증

```bash
go build ./...
go test ./... -count=1
filefunc validate .

# gozhip admin
tsma reset --all
tsma check
tsma status
# server.gen.go 함수가 목록에서 사라져야 함
# 오매칭 74건 해소되어야 함
```

## 의존성

없음.
