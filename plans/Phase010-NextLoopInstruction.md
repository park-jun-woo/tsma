# Phase010 — tsma next 출력에 다음 행동 지시 추가

## 목표

tsma next 출력 끝에 "작업 완료 후 `tsma next`를 실행하라"를 포함시켜, 모델이 임의 중단하지 않고 루프를 완주하도록 한다.

## 배경

- Codex, Grok 등 외부 모델이 tsma next를 4~5회 실행 후 임의 중단
- tsma next가 남았는데도 모델이 "충분히 했다"고 판단하고 멈춤
- 원인: tsma next 출력에 다음 행동이 명시되어 있지 않아 모델이 루프를 이어야 한다는 신호를 받지 못함

## 변경

tsma next 출력의 마지막에 다음 행동 지시를 추가:

### TODO 상태일 때 (테스트 작성 필요)

```
ListContracts  TODO
  file: internal/contract/service.go:41
  ▶ Write a test for this function.
  ▶ After completing the test, run `tsma next`.
```

### 커버리지 미달일 때 (추가 테스트 필요)

```
ListContracts  testing...
  [1/2] go test: PASS
  [2/2] coverage: 65% (11/17)
  UNCOVERED:
    line 41 — if params.Status != nil
    line 44 — if params.BuildingId != nil
  ▶ Cover the uncovered branches.
  ▶ After completing the test, run `tsma next`.
```

### PASS/DONE일 때 (다음 함수로 이동)

```
ListContracts  PASS (100%)
  ▶ Run `tsma next` for the next function.
```

### 전부 완료일 때

```
All functions complete!
```

이 경우에만 지시 없음 — 모델이 멈추는 게 맞음.

## 변경 파일

tsma next 출력을 생성하는 Go 코드. 정확한 파일은 구현 시 확인.

## 검증 방법

1. tsma next 출력 끝에 "After completing the test, run `tsma next`." 포함 확인
2. "All functions complete!" 시에만 지시 없음 확인
