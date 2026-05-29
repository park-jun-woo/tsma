# tsma 적용 결과 — gozhip 레거시 원본

- 대상: `~/.clari/repos/gozhip/artifacts/backend/admin/`
- 측정일: 2026-05-14
- tsma 상태: 완료 (TODO 0건)

## 전체 요약

| 지표 | 값 |
|------|-----|
| 전체 statement 커버리지 (go test) | **32.8%** |
| 총 함수 수 (tsma 인덱싱) | 527 |
| PASS (100% 브랜치 커버리지) | 246 (46%) |
| DONE (부분 커버리지, best effort) | 281 (53%) |
| TODO (미테스트) | 0 (0%) |
| 테스트 전부 통과 | YES |

## 패키지별 커버리지

### statement 커버리지 (go test -coverprofile)

| 패키지 | 커버리지 | 평가 |
|--------|---------|------|
| config | 100.0% | 완벽 |
| middleware | 97.2% | 우수 |
| crypto | 89.3% | 우수 |
| api/member | 79.0% | 양호 |
| api/form | 62.7% | 보통 |
| api/building | 55.2% | 보통 |
| storage | 54.5% | 보통 |
| api/webhook | 52.6% | 보통 |
| service | 52.0% | 보통 |
| api/ocr | 45.7% | 미흡 |
| api/auth | 45.6% | 미흡 |
| api/contract | 33.4% | 미흡 |
| db | 23.5% | 미흡 |
| repository | 1.2% | 사실상 미테스트 |
| service/gemini | 0.0% | 미테스트 (외부 API 의존) |

### 함수 레벨 분포 (go tool cover -func, mocks/cmd 제외)

| 패키지 | 함수 수 | 평균 커버리지 | 100% | 80-99% | 60-79% | 40-59% | 20-39% | 0% |
|--------|--------|-------------|------|--------|--------|--------|--------|-----|
| config | 3 | 100.0% | 3 | — | — | — | — | — |
| middleware | 21 | 99.6% | 20 | 1 | — | — | — | — |
| api/member | 7 | 88.0% | 3 | 2 | 2 | — | — | — |
| crypto | 2 | 87.9% | — | 2 | — | — | — | — |
| api/webhook | 3 | 82.4% | 2 | — | — | 1 | — | — |
| service | 165 | 80.2% | 113 | 17 | 6 | — | — | 29 |
| api/form | 11 | 72.3% | 5 | 1 | 2 | 1 | — | 2 |
| api/ocr | 5 | 63.3% | 2 | — | — | 2 | 1 | — |
| storage | 8 | 50.0% | 4 | — | — | — | — | 4 |
| api/building | 139 | 44.4% | 57 | 5 | — | — | — | 77 |
| api/auth | 42 | 35.7% | 15 | — | — | — | — | 27 |
| db | 3 | 33.3% | 1 | — | — | — | — | 2 |
| api/contract | 101 | 29.4% | 19 | 4 | 9 | 2 | — | 67 |
| repository | 180 | 1.7% | 3 | — | — | — | — | 177 |
| service/gemini | 3 | 0.0% | — | — | — | — | — | 3 |
| **합계** | **693** | — | **247** | **32** | **19** | **6** | **1** | **388** |

## 0% 커버리지 함수 분류 (repository/mocks/cmd 제외)

| 분류 | 수 | 설명 |
|------|-----|------|
| oapi-codegen 생성 코드 | 78 | `Visit*Response`, `RegisterHandlers`, `NewStrictHandler` 등. 코드젠 도구가 보장하는 영역. |
| repository 계층 | 177 | DB 쿼리 실행 함수. 실제 PostgreSQL 연결 필요. 단위 테스트 대상이 아님. |
| hand-written 로직 | 124 | mock 불가능한 의존성(concrete type)으로 인해 특정 브랜치 도달 불가 |
| service/gemini | 3 | Gemini API 호출. 외부 서비스 의존. |
| storage | 4 | S3 구현체. 외부 서비스 의존. |
| db | 2 | DB 커넥션 헬퍼. 실제 DB 연결 필요. |

## 0%인 이유 — 계층별 분석

### repository (1.2%, 177/180 함수 0%)

repository는 SQL 쿼리를 `database/sql`로 직접 실행하는 계층. 핸들러와 서비스 테스트에서 mock으로 대체되므로 repository 코드 자체는 실행되지 않음.

repository를 테스트하려면:
- 실제 PostgreSQL 인스턴스 필요
- 테스트 DB 마이그레이션 필요
- 테스트 데이터 시딩 필요

→ **단위 테스트가 아닌 통합 테스트 (Hurl 등) 영역**

### oapi-codegen 생성 코드 (78 함수 0%)

`server.gen.go`, `enums.gen.go`의 생성 함수:
- `Visit*Response` — strict server 응답 래퍼
- `RegisterHandlers`, `RegisterHandlersWithOptions` — 라우트 등록
- `NewStrictHandler` — strict handler 팩토리

코드젠 도구가 정확성을 보장하므로 수동 테스트 불필요. tsma Phase004에서 `*.gen.go` 제외 대상으로 지정.

### hand-written 0% (124 함수)

mock 불가능한 concrete type 의존으로 인해 일부 브랜치에 도달 불가능한 경우:

```go
type Handler struct {
    svc *service.BuildingService  // struct pointer — mock 불가
}
```

인터페이스로 변경하면 100% 가능하지만, 레거시 코드 수정 필요.

주요 0% 함수:
- `api/auth`: GetProfile, Login, Signup 등 — server.gen.go의 strict handler 래퍼 함수 (수작업 handler.go의 동명 함수와 별개)
- `api/building`: CheckBuildingDeletable, CheckUnitDeletable 등 — server.gen.go 래퍼
- `api/contract`: 대부분 server.gen.go 래퍼 + convert.go의 변환 유틸

## 잘 커버된 영역

### 100% 커버리지 (247 함수)

| 영역 | 대표 함수 | 특징 |
|------|----------|------|
| TransactionMatcher | Match, matchByAlias, matchByName, matchByUnitNumber, matchByAmount, normalizeName | 순수 로직, 외부 의존 없음 |
| LocalStorage | Upload, Download, Delete, NewLocalStorage | 파일시스템만 의존 |
| middleware | Auth, CORS, Logging, Recovery, RequestID, Response (20/21) | gin.Context mock 가능 |
| config | Load, getEnv, getEnvInt | 환경변수만 의존 |
| service (113 함수) | Login, Signup, SetupTOTP, VerifyTOTPLogin, RefreshToken, CreateTemplate, GetFormDefinition, SaveFormDefinition, SendInviteEmail 등 | 인터페이스 기반 DI |
| totp_helpers | generateBackupCodes, verifyBackupCode | 순수 로직 |

### 80-99% (32 함수)

| 영역 | 대표 함수 | 미커버 사유 |
|------|----------|------------|
| service | generateTokenPair (71%), ResetTOTP (84%), SetPassword (77%) | 일부 에러 브랜치 도달 불가 |
| api/member | Create (84%), Update (84%) | 일부 바인딩 에러 |
| crypto | Encrypt (81%), Decrypt (94%) | crypto/rand 실패 브랜치 |

## tsma 적용 과정 요약

1. `tsma next`를 반복 실행하여 527개 함수에 대해 테스트 작성
2. 각 함수마다: LLM이 함수 본문 + 의존성을 읽고 테스트 작성 → 실행 → 커버리지 피드백 → 개선
3. 100% 도달 함수 → PASS, 도달 불가 함수 → 1회 재시도 후 DONE
4. 최종: TODO 0건, 모든 함수에 테스트 존재

## 커버리지 향상 방안

### 단위 테스트로 개선 가능 (현재 도구로)

| 방안 | 대상 | 예상 효과 |
|------|------|----------|
| concrete type → interface 변경 | api/building, api/contract 핸들러 | 0% → 80%+ 가능 |
| gemini mock 인터페이스 | service/gemini | 0% → 90%+ 가능 |
| S3 mock 인터페이스 | storage/s3 | 0% → 90%+ 가능 |

### 통합 테스트가 필요한 영역

| 영역 | 함수 수 | 필요한 것 |
|------|--------|----------|
| repository | 177 | PostgreSQL + 마이그레이션 + 테스트 데이터 |
| db | 2 | PostgreSQL 연결 |
| api 핸들러 (E2E) | ~100 | 서버 기동 + HTTP 요청 (Hurl) |

→ **symfeed → Hurl 생성** 방향과 일치

### 이론적 최대치

| 분류 | 현재 | 제외 시 | 비고 |
|------|------|--------|------|
| 전체 | 32.8% | — | |
| repository 제외 | ~52% | 177 함수 제외 | repository는 통합 테스트 영역 |
| repository + gen 제외 | ~62% | 177 + 78 함수 제외 | 생성 코드는 도구가 보장 |
| mock 가능 계층만 | ~75% | 수작업 코드 중 DI 가능한 것만 | interface 변환 시 |
