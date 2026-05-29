# Phase011 — .tsmignore 파일 기반 인덱싱 제외

## 변경이력

- 2026-05-25: 작성. filefunc의 .ffignore 구현 참고.

## 목표

프로젝트 루트의 `.tsmignore` 파일에 기재된 패턴과 매칭되는 경로를 인덱싱에서 제외. filefunc의 `.ffignore`와 동일한 문법과 동작.

## 배경

tsma는 프로젝트의 모든 함수를 인덱싱하여 테스트 커버리지를 추적한다. 하지만 생성된 코드, 외부 의존성, 실험용 코드 등 테스트 대상이 아닌 파일이 인덱싱되면 TODO가 불필요하게 늘어난다.

현재 제외 수단은 `skipGoDir`/`skipPyDir`/`skipTSDir`에 하드코딩된 디렉토리명(`vendor`, `.git`, `.tsma`, `node_modules`)뿐. 프로젝트별 커스텀 제외가 불가능.

filefunc의 `.ffignore`는 이 문제를 해결하고 있으며, 동일한 구현을 tsma에 적용.

## filefunc 구현 분석 (확인 완료)

### 파일 3개로 구성

1. **`parse_ffignore.go`** — `.ffignore` 파일을 읽어 패턴 목록 반환. 빈 줄과 `#` 주석 skip. 파일 없으면 빈 슬라이스.
2. **`match_pattern.go`** — 단일 패턴과 경로 매칭. 3가지 분기:
   - `pattern/` (trailing slash) → 디렉토리 매칭
   - `path/to/pattern` (slash 포함) → 전체 경로 glob 매칭
   - `*.ext` (slash 없음) → 파일명 glob 매칭
3. **`match_ffignore.go`** — 패턴 목록 순회, 하나라도 매칭되면 true.

### 사용 흐름

```
ParseFFIgnore(root + "/.ffignore") → []string
  ↓
WalkFiles(root, ext, ignorePatterns) 내부에서:
  MatchFFIgnore(path, name, isDir, patterns) → bool
    → true + isDir → filepath.SkipDir
    → true + !isDir → skip file
```

## 변경 파일

| 파일 | 변경 |
|---|---|
| `internal/index/parse_tsmignore.go` (신규) | `.tsmignore` 파일 파싱. filefunc `ParseFFIgnore`와 동일 로직 |
| `internal/index/match_tsmignore.go` (신규) | 패턴 목록 매칭. filefunc `MatchFFIgnore`와 동일 로직 |
| `internal/index/match_pattern.go` (신규) | 단일 패턴 매칭. filefunc `matchPattern`과 동일 로직 |
| `internal/index/go_indexer_index.go` | `Index` 메서드에 `.tsmignore` 파싱 + walk 시 매칭 체크 추가 |
| `internal/index/py_indexer_index.go` | 동일 적용 |
| `internal/index/ts_indexer_index.go` | 동일 적용 |

## 변경 내용

### 1. parse_tsmignore.go

```go
func ParseTsmIgnore(path string) []string {
    // .ffignore와 동일: 파일 읽기, 빈 줄/#주석 skip, 패턴 목록 반환
}
```

### 2. match_pattern.go + match_tsmignore.go

```go
func matchPattern(path, name string, isDir bool, pattern string) bool {
    // filefunc와 동일: trailing slash → dir, slash 포함 → path glob, 그 외 → name glob
}

func MatchTsmIgnore(path, name string, isDir bool, patterns []string) bool {
    // 패턴 목록 순회, 하나라도 매칭 → true
}
```

### 3. 각 Indexer.Index() 수정

```go
func (g *GoIndexer) Index(projectRoot string) ([]model.Function, error) {
    ignorePatterns := ParseTsmIgnore(filepath.Join(projectRoot, ".tsmignore"))
    // ...
    err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
        // 기존 skipGoDir 이전에 .tsmignore 체크 추가
        if MatchTsmIgnore(path, info.Name(), info.IsDir(), ignorePatterns) {
            if info.IsDir() { return filepath.SkipDir }
            return nil
        }
        // ... 기존 로직 유지
    })
}
```

Py/TS indexer도 동일 패턴.

## .tsmignore 문법 (filefunc .ffignore와 동일)

```
# 주석
vendor/          # 디렉토리 제외 (trailing slash)
internal/db/     # 특정 하위 디렉토리 제외
*.gen.go         # 파일명 패턴 제외
cmd/main.go      # 특정 파일 경로 제외
```

## 의존성

없음. filefunc 코드를 참고하되 복사 구현 (tsma는 독립 모듈).

## 검증

```bash
cd ~/.clari/repos/fullend/tsma && go test ./internal/index/... -v
cd ~/.clari/repos/fullend/tsma && go test ./...
```

테스트 케이스:
- .tsmignore 없으면 전체 인덱싱 (기존 동작 보존)
- 디렉토리 패턴 (`vendor/`) → 해당 디렉토리 전체 skip
- 파일 패턴 (`*.gen.go`) → 매칭 파일 skip
- 경로 패턴 (`internal/db/models.go`) → 특정 파일 skip
- 주석/빈 줄 → 무시
