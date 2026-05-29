# Phase005 — Python/TypeScript 함수 EndLine 계산

## 문제

Python·TypeScript 인덱서 모두 `EndLine = StartLine`으로 설정한다. 함수 범위가 1줄이 되어 커버리지 체커가 선언 라인 하나만 검사하고 본문 전체를 무시한다. 사실상 per-function 커버리지 분석이 무의미.

- Python: `match_py_func.go:45-46`
- TypeScript: `match_ts_top_level_func.go:34`, `match_ts_method.go:36`

Go 인덱서는 `go/ast`가 `FuncDecl.End()` 위치를 제공하므로 이 문제가 없다.

## 수정 방향

### Python

들여쓰기 기반으로 함수 본문 끝을 추적한다. `def`/`async def` 라인의 들여쓰기보다 깊은 라인이 연속되는 동안 EndLine을 갱신. 빈 줄·주석만 있는 줄은 건너뛴다. 다음 함수/클래스 선언 또는 동일·상위 들여쓰기의 실행문을 만나면 종료.

```
def login(email, password):   ← StartLine=10
    user = find(email)         ← 본문 (indent > def indent)
    if not user:
        raise NotFound()
    return user                ← EndLine=14
                               ← 빈 줄, 건너뜀
def logout():                  ← 새 함수, login 종료 확정
```

### TypeScript

중괄호 깊이 카운터로 함수 본문 끝을 추적한다. 함수 선언 라인에서 `{` 를 만나면 depth=1 시작, `}` 로 depth=0이 되면 EndLine 확정.

문자열·주석 내 중괄호를 무시하기 위한 간이 스캐너 규칙:

| 컨텍스트 | 진입 | 탈출 | 중괄호 카운팅 |
|---------|------|------|-------------|
| 한 줄 주석 | `//` | 줄 끝 | 무시 |
| 여러 줄 주석 | `/*` | `*/` | 무시 |
| 작은따옴표 문자열 | `'` | `'` (이스케이프 `\'` 제외) | 무시 |
| 큰따옴표 문자열 | `"` | `"` (이스케이프 `\"` 제외) | 무시 |
| 템플릿 리터럴 | `` ` `` | `` ` `` (이스케이프 제외) | `${` 내부만 카운팅 (중첩 depth 별도 추적) |

정규식 리터럴(`/pattern/`) 내 중괄호는 처리하지 않는다. 정규식과 나눗셈을 문맥 없이 구분하는 것은 파서 수준 작업이며, 실전에서 정규식 내 `{`/`}`가 EndLine 오판을 일으킬 확률은 극히 낮다.

```
export function login(req: Request): Response {   ← StartLine=10, depth=1
  const user = findUser(req.email);                ← depth=1
  if (!user) {                                     ← depth=2
    throw new NotFoundError();
  }                                                ← depth=1
  return user;
}                                                  ← depth=0, EndLine=16
```

#### arrow function 처리

중괄호가 있는 경우 (`const fn = () => { ... }`) — 동일 로직 적용.

중괄호가 없는 경우 (`const fn = () => expr`) — 선언 라인에서 `{`가 나타나지 않으면, 다음 조건 중 하나를 만날 때까지 EndLine 갱신:
- 동일 들여쓰기 이하의 비공백 라인 (새로운 선언)
- 세미콜론으로 끝나는 라인
- 빈 줄 다음에 동일 들여쓰기 이하의 라인

이 휴리스틱은 100% 정확하지 않으나, 중괄호 없는 arrow function이 여러 줄에 걸치는 경우 자체가 드물다.

## 변경 파일

| 파일 | 변경 |
|---|---|
| `internal/index/index_py_file.go` | 함수 감지 후 후속 라인 스캔으로 EndLine 계산 |
| `internal/index/match_py_func.go` | EndLine 파라미터 제거, index_py_file에서 직접 설정 |
| `internal/index/index_ts_file.go` | 함수 감지 후 중괄호 깊이 추적으로 EndLine 계산 |
| `internal/index/match_ts_top_level_func.go` | EndLine 파라미터 제거 |
| `internal/index/match_ts_method.go` | EndLine 파라미터 제거 |

## 검증

```bash
go test ./internal/index/... -count=1 -run "Py\|TS"
go test ./internal/coverage/... -count=1

# 더미 프로젝트로 EndLine 확인
tsma reset --all && tsma check
tsma list  # EndLine이 StartLine보다 큰지 육안 확인
```

## 의존성

없음. Phase004(파일 기반 매칭)와 독립.
