# Phase001 — string concatenation 라우트 경로 탐지

## 목표

`router.GET(options.BaseURL+"/api/v1/buildings", wrapper.ListBuildings)` 같은 `*ast.BinaryExpr` 패턴의 라우트 경로를 추출하여, oapi-codegen 등 코드 생성기가 만든 엔드포인트를 누락 없이 탐지한다.

## 배경

gozhip admin 프로젝트 실사용 테스트에서 108개 라우트 중 44개가 누락됨. 원인: `server.gen.go`(oapi-codegen 생성)의 라우트 등록이 string literal이 아닌 concatenation expression 사용.

```go
// 현재 탐지 가능 (BasicLit)
r.GET("/api/v1/auth/login", h.Login)

// 현재 누락 (BinaryExpr)
router.GET(options.BaseURL+"/api/v1/admin/buildings", wrapper.ListBuildings)
```

## 영향 범위

Go 3개 detector 모두 동일 패턴으로 경로를 추출하므로 공통 함수로 해결.

| 파일 | 현재 코드 |
|---|---|
| `internal/endpoint/match_gin_route.go:28-31` | `call.Args[0].(*ast.BasicLit)` |
| `internal/endpoint/match_echo_route.go:34-37` | 동일 |
| `internal/endpoint/match_chi_route.go:30-33` | 동일 |
| `internal/endpoint/match_echo_add_route.go:21-24` | 동일 (두 번째 인자) |

## 변경 계획

### 1. 새 파일: `internal/endpoint/extract_route_path.go`

`ast.Expr`에서 라우트 경로를 추출하는 공통 함수.

```go
func extractRoutePath(expr ast.Expr) (string, bool)
```

처리 순서:
1. `*ast.BasicLit` (STRING) → 그대로 반환. 기존 동작 보존.
2. `*ast.BinaryExpr` (ADD) → 좌항·우항을 재귀 순회하여 string literal 부분만 연결. 변수/함수 호출 부분은 건너뜀.
   - `options.BaseURL + "/api/v1/buildings"` → `"/api/v1/buildings"`
   - `prefix + "/users/" + suffix` → `"/users/"`
3. 그 외 → `("", false)`

### 2. 기존 4개 match 파일 수정

`call.Args[0].(*ast.BasicLit)` 패턴을 `extractRoutePath(call.Args[0])` 호출로 교체.

| 파일 | 변경 |
|---|---|
| `match_gin_route.go` | L28-32 → `extractRoutePath` 호출 |
| `match_echo_route.go` | L34-38 → 동일 |
| `match_chi_route.go` | L30-34 → 동일 |
| `match_echo_add_route.go` | L21-25 → 동일 (두 번째 인자에 적용) |

### 3. 테스트: `internal/endpoint/extract_route_path_test.go`

| 케이스 | 입력 | 기대 결과 |
|---|---|---|
| 순수 string literal | `"/api/v1/users"` | `"/api/v1/users"` |
| 변수 + literal | `options.BaseURL + "/api/v1/buildings"` | `"/api/v1/buildings"` |
| literal + 변수 | `"/api/v1/" + pathSuffix` | `"/api/v1/"` |
| 변수 + literal + 변수 | `prefix + "/users/" + suffix` | `"/users/"` |
| 순수 변수 | `dynamicPath` | `("", false)` |
| 다단 concatenation | `a + "/b" + "/c"` | `"/b/c"` |

### 4. 통합 테스트: Gin detector가 BinaryExpr 패턴 탐지 확인

기존 `go_gin_test.go` 패턴에 추가 — fixture에 concatenation 라우트를 포함하여 탐지 여부 확인.

## 검증

```bash
go test ./internal/endpoint/ -cover -v
filefunc validate /path/to/testmaster
```

- 기존 테스트 전통과 (회귀 없음)
- 새 테스트 통과
- gozhip admin에서 `testmaster reset --all && testmaster next` → 108개 탐지

## 의존성

없음. endpoint 패키지 내부 변경.
