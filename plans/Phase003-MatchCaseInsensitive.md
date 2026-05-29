# Phase003 — match 대소문자 불일치 수정

## 버그

비공개 함수의 테스트가 존재하는데 TODO로 남는다.

- `func getEnv()` → `func TestGetEnv()` — 매칭 실패
- `func parseTemplateID()` → `func TestParseTemplateID()` — 매칭 실패
- `func generateID()` → `func TestGenerateID()` — 매칭 실패
- `func toBuildingResponse()` → `func TestToBuildingResponse()` — 매칭 실패

## 근본 원인

`internal/match/matches_test_func.go:25`:

```go
suffix := testName[len("Test"):]
return strings.Contains(suffix, funcName)
```

Go 관례: 비공개 함수 `getEnv`의 테스트는 `TestGetEnv` (Test + 첫 글자 대문자). `strings.Contains("GetEnv_Fallback", "getEnv")` → false. 대소문자가 다르면 매칭 실패.

## 수정

`strings.Contains` → `strings.EqualFold` 기반 비교, 또는 첫 글자를 대문자로 변환 후 비교.

Go 관례에 맞는 정확한 방법: funcName의 첫 글자를 대문자로 올린 뒤 Contains.

```go
capitalized := strings.ToUpper(funcName[:1]) + funcName[1:]
return strings.Contains(suffix, capitalized) || strings.Contains(suffix, funcName)
```

양쪽 다 체크하여 exported/unexported 함수 모두 매칭.

## 변경 파일

| 파일 | 변경 |
|---|---|
| `internal/match/matches_test_func.go:25` | capitalize 후 비교 추가 |

## 검증

```bash
# gozhip admin에서 재확인
tsma reset --all
tsma next
tsma check
tsma status
# getEnv, parseTemplateID, generateID, toBuildingResponse가 DONE이어야 함
```

## 의존성

없음. match 패키지 1줄 수정.
