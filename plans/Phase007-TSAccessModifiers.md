# Phase007 — TypeScript 접근 제어자·const 비함수 필터링

## 문제

### 1. 접근 제어자가 함수명으로 잡힘

메서드 패턴 (`ts_patterns.go:17`)이 `private`/`public`/`protected`/`static`/`override`/`readonly`/`abstract` 키워드를 고려하지 않음.

```typescript
class AuthService {
  private async validate(data: any) {  // "private"가 함수명으로 캡처됨
```

### 2. `const` 비함수 할당이 함수로 인덱싱

`tsFuncPattern`의 두 번째 대안 (`^(?:export\s+)?(?:const|let|var)\s+(\w+)\s*=`)이 화살표 함수뿐 아니라 모든 할당을 매칭.

```typescript
export const MAX_RETRIES = 5;       // 함수로 인덱싱됨 (false positive)
export const API_URL = "/api/v1";   // 함수로 인덱싱됨 (false positive)
```

## 수정 방향

### 접근 제어자

메서드 패턴 앞에 선택적 접근 제어자 그룹 추가:

```
// 현재
^\s+(?:async\s+)?(\w+)\s*\([^)]*\)\s*[:{]
// 변경
^\s+(?:(?:private|public|protected|static|override|readonly|abstract)\s+)*(?:async\s+)?(\w+)\s*(?:<[^>]*>)?\s*\([^)]*\)\s*[:{]
```

제네릭 메서드 (`method<T>(arg: T)`)도 `(?:<[^>]*>)?`로 처리.

### const 비함수 필터링

라인에 `=>` 또는 `function`이 포함되어 있는지 추가 확인. 없으면 스킵.

```go
// matchTSTopLevelFunc 내부
if submatch[2] != "" {  // const/let/var 매칭
    if !strings.Contains(line, "=>") && !strings.Contains(line, "function") {
        return nil  // 화살표 함수도 function도 아님 → 스킵
    }
}
```

단, 여러 줄에 걸친 화살표 함수 (`const fn = (\n) => {`)는 선언 줄에 `=>`가 없을 수 있음. 이 경우는 미감지 허용 — regex 기반 인덱서의 구조적 한계이며, 전체 파서(tree-sitter 등)를 도입하지 않는 한 해결 불가. 실전에서 여러 줄 화살표 함수 선언은 드물고, 미감지 시 해당 함수의 테스트가 TODO로 남지 않을 뿐 다른 함수에 영향을 주지 않는다.

## 변경 파일

| 파일 | 변경 |
|---|---|
| `internal/index/ts_patterns.go` | 메서드 패턴에 접근 제어자 + 제네릭 추가 |
| `internal/index/match_ts_top_level_func.go` | const 매칭 시 `=>` / `function` 존재 확인 |
| `internal/index/match_ts_method.go` | 스킵 리스트에 `private`, `public`, `protected`, `static`, `override`, `readonly`, `abstract` 불필요 (패턴이 소비하므로) |

## 검증

```bash
go test ./internal/index/... -count=1 -run "TS"

# 테스트 케이스 추가
# - private async validate(...) → 함수명 "validate"
# - public static create(...) → 함수명 "create"
# - const MAX = 5 → 인덱싱 안 됨
# - const handler = async () => { → 인덱싱 됨
```

## 의존성

Phase006 선행 (`.tsx` 지원 후 접근 제어자 테스트 가능). 또는 `.ts` 파일만으로 독립 검증 가능.
