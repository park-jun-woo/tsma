# Phase006 — TypeScript `.tsx`/`.jsx` 지원 + `export default` 패턴

## 문제

### 1. `.tsx`/`.jsx` 파일 완전 무시

`isTSSource()` (`is_ts_source.go:18`)가 `.ts`와 `.js`만 허용. React 프로젝트의 모든 컴포넌트·훅·페이지가 `.tsx`이므로 인덱싱 0건.

### 2. `.test.tsx`/`.spec.tsx` 테스트 미매칭

`tsTestSuffixes` (`ts_test_suffixes.go:4`)가 `.test.ts`, `.test.js`, `.spec.ts`, `.spec.js`만 포함. React 테스트 파일(`Component.test.tsx`)을 찾지 못함.

### 3. `export default function`/`export default class` 미매칭

`tsFuncPattern`, `tsClassPattern` (`ts_patterns.go:10,15`)이 `default` 키워드를 소비하지 않음. Next.js 페이지 (`export default function Page()`) 등 감지 불가.

## 수정 방향

### `.tsx`/`.jsx` 인덱싱

`isTSSource()`에 `.tsx`, `.jsx` 확장자 추가. 테스트 제외 패턴도 `.test.tsx`, `.test.jsx`, `.spec.tsx`, `.spec.jsx` 추가.

### 테스트 매칭

`tsTestSuffixes`에 `.test.tsx`, `.test.jsx`, `.spec.tsx`, `.spec.jsx` 4개 추가.

### `export default` 패턴

정규식에 `(?:default\s+)?` 그룹 삽입:

```
// 현재
^(?:export\s+)?(?:async\s+)?function\s+(\w+)
// 변경
^(?:export\s+)?(?:default\s+)?(?:async\s+)?function\s+(\w+)

// 현재
^(?:export\s+)?class\s+(\w+)
// 변경
^(?:export\s+)?(?:default\s+)?class\s+(\w+)
```

## 변경 파일

| 파일 | 변경 |
|---|---|
| `internal/index/is_ts_source.go` | `.tsx`, `.jsx` 허용 + `.test.tsx` 등 테스트 제외 |
| `internal/index/ts_patterns.go` | `(?:default\s+)?` 그룹 추가 (함수, 클래스) |
| `internal/match/ts_test_suffixes.go` | `.test.tsx`, `.test.jsx`, `.spec.tsx`, `.spec.jsx` 추가 |

## 검증

```bash
go test ./internal/index/... -count=1 -run "TS"
go test ./internal/match/... -count=1 -run "TS"

# React 프로젝트(gozhip arts/admin/frontend)로 실제 검증
cd ~/.clari/repos/gozhip/arts/admin/frontend
tsma reset --all && tsma check
tsma status  # .tsx 파일의 함수가 인덱싱되었는지 확인
```

## 의존성

없음. Phase005와 독립 (병렬 가능).
