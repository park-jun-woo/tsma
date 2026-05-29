# Phase009 — Python/TypeScript 통합 테스트 더미 프로젝트

## 문제

Python·TypeScript 지원의 전체 파이프라인(index → match → run → coverage)이 실제 프로젝트에서 end-to-end로 검증된 적 없음. 모든 테스트가 Go 단위 테스트(임시 디렉토리 + writeFile)로만 수행됨.

Phase005~008의 수정이 실제로 동작하는지 확인할 수단이 필요.

## 수정 방향

`dummys/tsma/` 아래에 최소 더미 프로젝트 2개 생성. 각 프로젝트는 tsma의 전체 워크플로(index → match → run → coverage)를 검증할 수 있는 최소 구조.

### Python 더미: `dummys/tsma/pyapp/`

```
pyapp/
├── app/
│   ├── __init__.py
│   ├── auth.py          # login(), signup() — 분기 포함
│   ├── building.py      # create_building(), list_buildings()
│   └── utils.py         # format_currency() — 단순 함수
├── tests/
│   ├── __init__.py
│   ├── test_auth.py     # login 테스트만 (signup은 TODO로 남김)
│   └── test_utils.py    # format_currency 100% 커버
├── requirements.txt     # pytest, pytest-cov
└── pyproject.toml       # [tool.pytest.ini_options]
```

검증 항목:
- `tsma check` → auth.py (login DONE, signup TODO), utils.py (DONE), building.py (TODO)
- `tsma next` → signup 또는 building 함수로 안내, 테스트 실행 + 커버리지 측정
- EndLine이 정확히 함수 끝을 가리키는지 (Phase005)
- pytest 기반 커버리지가 정상 동작하는지
- `tests/` 디렉토리 매칭이 동작하는지

### TypeScript 더미: `dummys/tsma/tsapp/`

```
tsapp/
├── src/
│   ├── auth.ts          # export function login(), export default function signup()
│   ├── building.tsx     # export default function BuildingList() — React 컴포넌트
│   ├── service.ts       # export class AuthService { private validate(), public login() }
│   └── utils.ts         # export const formatCurrency = () => {}, const MAX = 5
├── src/__tests__/
│   ├── auth.test.ts     # login 테스트만
│   └── utils.test.ts    # formatCurrency 100% 커버
├── package.json         # vitest, @vitest/coverage-v8 in devDependencies
├── vitest.config.ts
└── tsconfig.json
```

검증 항목:
- `.tsx` 파일 인덱싱 (Phase006)
- `export default function` 인덱싱 (Phase006)
- `private`/`public` 메서드 정확한 이름 추출 (Phase007)
- `const MAX = 5` 미인덱싱, `const formatCurrency = () =>` 인덱싱 (Phase007)
- `__tests__/` 매칭, `.test.ts` 매칭
- vitest 커버리지 동작
- EndLine 정확성 (Phase005)

### 환경 세팅

더미 프로젝트는 의존성을 미리 설치해두지 않는다. 검증 스크립트가 설치를 수행.

### 검증 스크립트

`dummys/tsma/run_integration.sh`:

```bash
#!/bin/bash
set -e

echo "=== Python integration test ==="
cd "$(dirname "$0")/pyapp"
python3 -m pip install -r requirements.txt --quiet
tsma reset --all
tsma check
tsma status
# 기대: 5 functions, 2 DONE, 3 TODO
tsma next
# 기대: TODO 함수 안내 → 테스트 실행은 하지 않음 (LLM 작성 단계)
# 러너·커버리지 검증: 기존 테스트 파일이 있는 함수에 대해 실행 + 커버리지 측정이 동작하는지 확인

echo ""
echo "=== TypeScript integration test ==="
cd "$(dirname "$0")/tsapp"
npm install --silent
tsma reset --all
tsma check
tsma status
# 기대: .tsx 포함 인덱싱, const MAX 미포함, private 정확 파싱
tsma next
# 기대: TODO 함수 안내 → vitest 실행 + 커버리지 측정이 동작하는지 확인
```

## 변경 파일

| 파일 | 변경 |
|---|---|
| `dummys/tsma/pyapp/` | Python 더미 프로젝트 신규 |
| `dummys/tsma/tsapp/` | TypeScript 더미 프로젝트 신규 |
| `dummys/tsma/run_integration.sh` | 통합 테스트 스크립트 (환경 세팅 포함) |

## 검증

```bash
cd dummys/tsma && bash run_integration.sh
```

수동 확인 사항:
1. Python: `tsma list`에서 각 함수의 EndLine이 StartLine보다 큰지
2. TypeScript: `tsma list`에서 `const MAX`가 목록에 없는지, `private validate`가 `validate`로 잡혔는지
3. 양쪽 모두: `tsma next`에서 기존 테스트 있는 함수의 커버리지가 0%가 아닌 값으로 측정되는지

## 의존성

Phase005(EndLine), Phase006(TSX), Phase007(접근 제어자), Phase008(Python 러너) 모두 완료 후 실행. 이 Phase는 최종 통합 검증.
