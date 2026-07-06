# tsma rulebook — 레거시 게이트 규칙 정리 (reins 포팅 기준)

> 목적: tsma가 **암묵적으로 손구현**해 온 결정론적 게이트(판정) 규칙을 한곳에 모아,
> reins `gate.Definition`/`gate.Rule`로 재표현할 수 있는 형태로 정리한다.
> tsma는 reins 비의존(현재 cobra만 의존)이며, 아래 규칙들은 `internal/cli`·`internal/coverage`·
> `internal/model`·`internal/match`·`internal/index`에 흩어져 있는 판정 로직을 추출한 것이다.
>
> 표기: 각 규칙에 임시 ID를 부여한다. `G-*`=게이트(verdict), `S-*`=상태/래칫, `E-*`=Item 적격성,
> `M-*`=매칭, `TS-*`=test-smell(미구현 제안). 출처는 `file:symbol`로 명시한다.

---

## 0. 한눈에 보는 reins 매핑

| reins 개념 | tsma 현재 손구현 | 비고 |
|---|---|---|
| `quest.Item` | `model.Function` (`function.go`) | Key=QualifiedName, State=Status, Tries=Attempt |
| `quest.Session` | `model.Session` (`session.go`) | Functions/Summary/CurrentIndex/MaxAttempts |
| `quest.State` | `Status` todo/pass/done (`status_constants.go`) | REVIEW/SKIPPED/BLOCKED **없음** |
| `quest.Apply` (PASS lock / FAIL Tries++ / MaxTries→DONE) | `apply_pass/done/retry_result.go` + `attemptOutcome` | 흩어져 있음 |
| `gate.Evaluate` → Verdict | `run_and_measure.go`의 outcome 분기 | 단일 지표(커버리지)라 룰 카탈로그 아님 |
| `gate.Rule` 카탈로그 / `rules` 명령 | **없음** | reins가 공급하는 핵심 |
| `Seed` (input→TODOs) | `analyze_project.go` 인덱싱 | E-* 규칙 |
| `Prepare` (제출 디코드/부작용) | `runAndMeasure`(테스트 실행+커버리지 측정) | 부작용은 여기로 |
| `Render` (작성 프롬프트) | `print_next_instruction.go` 등 | |

---

## 1. 상태 머신 / 래칫 (S-*)

tsma의 상태는 3개뿐이다. (`internal/model/status_constants.go`)

```
StatusTodo = "todo"   StatusPass = "pass"   StatusDone = "done"
```

### S-001 — 상태 집합과 종결성
- **todo**: 미완. 다음 라운드 대상.
- **pass**: 브랜치 커버리지 100% 달성. **종결·불가역**.
- **done**: 100% 미만이나 best-effort로 수용(자동/명시). **종결**.
- 근거: `session_recalc_summary.go` — pass/done만 카운트, 그 외 전부 todo로 집계(미지·"fail"·"" 포함).
- reins 차이: tsma엔 `REVIEW`/`SKIPPED`/`BLOCKED`가 없다. test-smell(TS-*)을 도입하면 `REVIEW`가 필요.

### S-002 — 래칫 단조성 (불가역)
- 한번 `pass`가 되면 되돌리지 않는다. 남은 작업(todo)은 단조 감소.
- 단, **stale 세션 거짓 완료 방지**(BUG-004, v0.4.1): `next`/`status`는 소스를 재스캔해
  실제 상태와 어긋난 종결을 막는다. → reins `quest.Apply`의 lock과 동형이되, 소스 변경 감지로 보정.

### S-003 — MaxAttempts 자동 DONE 강등 (`attempt_outcome.go`)
```
attemptOutcome(attempt, maxAttempts):
    attempt >= maxAttempts → outcomeDone   // best-effort 수용
    else                   → outcomeRetry  // todo 유지
```
- reins의 `MaxTries=3 → DONE`과 정확히 동형. tsma는 `Session.MaxAttempts`로 가변(기본값 `defaultMaxAttempts`, `--max-attempts` 플래그).
- **단, 자동 강등은 "변경 없는 재제시(resurface)" 경로에서만** 카운트된다 — 아래 S-004 참조.

### S-004 — 정지 방지 / 회전 커서 (BUG-002, `run_next_interactive.go`)
- 1차 패스 이후 `CurrentIndex`는 남은 TODO 위를 **회전**한다.
- **변경 없는 partial**(테스트 미수정): 한 번 보여주고 커서를 돌린다 →
  단일 partial이 뒤의 TODO들을 막지 못한다. 재제시 횟수가 `maxAttempts`에 닿으면 DONE 수용
  (`resurface_partial.go` → `attemptOutcome`).
- **실제 편집된 partial**(`outcomeRetry`): 커서를 그대로 두어 방금 만진 함수를 계속 보여준다.
  이 경로는 자동 DONE으로 강등하지 않는다(사용자가 작업 중이므로).

---

## 2. 핵심 게이트 — 커버리지 verdict (G-*)

verdict는 `run_and_measure.go`에서 결정된다. **단일 연속 지표(브랜치 커버리지)** 기반이라 reins
포팅 시 "룰 카탈로그"가 아니라 **단일 Fail 룰 + outcome 매핑**으로 표현된다.

판정 순서(우선순위 = 위에서 아래, 먼저 걸리는 것이 승):

### G-001 — 테스트 실행 실패 = TEST_FAIL (최우선)
- 매칭된 테스트를 실행(`runner.Run`)했을 때 `err != nil || !res.Pass` → `outcomeTestFail`.
- 커버리지 측정 단계(`checker.Check`)에서 에러가 나도 동일하게 TEST_FAIL.
- 의미: **컴파일/실행이 깨지면 커버리지는 보지 않는다.** todo 유지 + 실패 출력 노출.
- reins 매핑: `LevelFail` 룰 `tests-must-pass`. (가장 상류, supersedes 대상.)

### G-002 — 브랜치 100% = PASS
- `report.AllCovered == true` → `outcomePass`, coveragePct=100, **상태 pass 잠금**.
- `AllCovered`의 정의(`build_go_report.go`): **함수의 모든 커버 블록이 100%** 일 때만 true.
  함수 하나라도 `CoveredPct < 100`이면 false + 미커버 브랜치 누적.
- 다국어 동형: `build_{go,java,py,rs,ts}_report.go` 전부 `Report{AllCovered:true}`에서
  미커버 발견 시 false로 떨어뜨리는 동일 구조.
- reins 매핑: **중심 게이트 룰** `branch-coverage-below-100` (`LevelFail`). 미발화 = PASS.
  Fact = `{Where: 미커버 라인, Expected: "100% branch", Actual: TotalPct}`.

### G-003 — 100% 미만 + 시도 소진 = DONE (자동 수용)
- `AllCovered == false` 이고 `attempt(=fn.Attempt+1) >= maxAttempts` → `outcomeDone`.
- coveragePct = `report.TotalPct`(실측), 상태 done 잠금.
- reins 매핑: G-002 룰이 발화했으나 `Tries==MaxTries` → reins가 자동으로 DONE 잠금(프레임워크 기본 동작).

### G-004 — 100% 미만 + 시도 남음 = RETRY (todo 유지)
- `AllCovered == false` 이고 `attempt < maxAttempts` → `outcomeRetry`.
- 미커버 브랜치 목록(`report.Uncovered`)을 피드백으로 노출 → 다음 시도 유도.
- reins 매핑: `OutFail` + Fact 피드백 → `Tries++` → 재시도. **이것이 tsma의 핵심 수렴 루프.**
- 비고: README 검증 — 미커버 라인 번호 피드백을 주면 LLM이 한 방에 100% 도달(피드백 패리티 효과).

> **요약 게이트 결정표**
>
> | 조건 | outcome | 상태 | reins verdict |
> |---|---|---|---|
> | 테스트 실패/측정 에러 | test_fail | todo | FAIL(tests-must-pass) |
> | AllCovered | pass | **pass(lock)** | PASS |
> | !AllCovered & tries 소진 | done | **done(lock)** | FAIL→MaxTries DONE |
> | !AllCovered & tries 남음 | retry | todo | FAIL + Fact 피드백 |

---

## 3. Item 적격성 / 인덱싱 규칙 (E-*)

`Seed`에 해당. 무엇이 "테스트해야 할 함수(Item)"가 되는가. (`internal/index`)

### E-001 — Go 소스 적격성 (`is_go_source.go`)
다음은 **인덱싱 제외**(Item 아님):
- `*_test.go` (테스트 파일 자체)
- `*_gen.go` / `*.gen.go` / `*.pb.go` (생성물)
- `mock_*` 프리픽스 파일
- `.go` 확장자 아님
- 함의: **test-smell(TS-*) 룰은 `_test.go`가 표적인데 현재 인덱서는 이를 제외** → smell은 별도 walk 필요(rulebook §6).

### E-002 — .tsmignore 필터 (`match_tsmignore.go`, `parse_tsmignore.go`)
- 프로젝트 루트 `.tsmignore` 패턴에 걸리는 경로/디렉터리는 walk에서 제외(`SkipDir`).
- reins 포팅 시 재사용 가능(파일 적격 필터 E-001과 달리 이건 그대로 이식).

### E-003 — 디렉터리 스킵 (`skip_{go,cs,java,py,rs,ts}_dir.go`)
- vendor·빌드 산출물 등 언어별 표준 제외 디렉터리.

### E-004 — 다국어 인덱싱 방식 차이 (중요)
- **Go만 go/ast**(`go_indexer_index.go` `parser.ParseFile`). 정밀.
- C#/Java/Rust/TS/Py는 **라인 기반 렉시컬 파서**(`dispatch_*_line.go`, `count_braces.go`,
  `match_pattern.go`). → 정밀도 한 단계 낮음. 룰 정밀도 기대치를 언어별로 다르게 잡아야 함.

---

## 4. 매칭 게이트 — 테스트 귀속 (M-*)

함수 ↔ 테스트 연결. 잘못 매칭하면 게이트가 엉뚱한 테스트로 판정한다. (`internal/match`)
결과 타입 `TestMatch{Files, TestFuncs}` — `TestFuncs` 비면 "Files의 모든 테스트 실행"(`test_match.go`).

### M-001 — 호출 기반 매칭 (Go, 주 경로)
- 테스트가 **실제로 호출하는** 소스 심볼로 귀속(`go_func_matcher.go`, `collect_called_refs.go`).
- 리시버 구분(`callee_receiver.go`, `filter_refs_by_receiver.go`) — 메서드 오매칭 방지.

### M-002 — receiver 오매칭 방지 (BUG-003, v0.3.4)
- 같은 이름·다른 리시버 메서드를 구분(`is_same_name_multiple.go`, `keep_ref_for_receiver.go`).
- 이게 깨지면 A.Foo 테스트가 B.Foo로 귀속되는 거짓 PASS 발생.

### M-003 — 파일명 fallback (BUG-001, v0.3.3 복원)
- 호출 기반 매칭 실패 시 파일명 규칙으로 fallback(`go_filename_fallback.go`,
  `fallback_func_matcher.go`). 비-Go는 fallback이 주 경로(`TestFuncs` nil → 파일 전체 실행).

### M-004 — 오명명 테스트 탐지 (`find_misnamed.go`)
- 함수에 맞는 테스트가 "거의 맞는 이름"으로 존재하면 rename 힌트 출력(`FindMisnamedTest`).
- reins 매핑: `Prepare`/`Render` 단계의 안내. 게이트 verdict 자체는 아님.

### M-005 — 매칭 일관성 (`match_funcs.go`)
- 배치 매칭(1차 풀스캔)과 단건 매칭(`detectTestChange`)이 **동일 규칙**으로 매칭해야 함.
- 어긋나면 1차/증분 결과 불일치.

---

## 5. 진입/스캔 절차 규칙 (P-*)

### P-001 — 1차 풀스캔 후 회전 모드 (`analyze_project.go`, v0.4.0)
- 최초 `next`는 전 함수 풀스캔(배치 매칭+측정) → `FirstPassDone=true`.
- 이후 `next`는 회전 커서 증분 모드(`run_next_interactive.go`, §S-004).

### P-002 — 변경 감지 재측정 (`detectTestChange`)
- 테스트 파일 mtime(여러 파일이면 max, `combinedTestMtime`)이 기록값과 다르면 "changed" → 재측정.
- 변경 없으면 재측정 생략(resurface 경로). → 불필요한 테스트 재실행 회피.

### P-003 — stale 세션 재스캔 (BUG-004, v0.4.1)
- `next`/`status`가 소스를 재스캔해, 세션 캐시가 실제와 어긋난 거짓 완료를 보정(§S-002).

---

## 6. test-smell 게이트 — 미구현 제안 (TS-*)

`files/reflect-rule.md`의 제안. **커버리지(G-*)와 분리된 품질 축.** 압박받은 LLM이 커버리지를
cheese(`unsafe`/`reflect` 동적 침투/`linkname`으로 내부 강제 도달)하는 것을 막는다.
reins 포팅 시 **추가 `gate.Rule`로 자연스럽게 편입**되며 `LevelReview` 권장.

| ID | 탐지 | reins Level | 비고 |
|---|---|---|---|
| TS-REFL-001 | `import "unsafe"` 또는 `unsafe.` 셀렉터 | Review(또는 Fail) | 강한 cheese |
| TS-REFL-002 | `reflect` 값의 `.MethodByName(`/`.FieldByName(` | Review | 비공개 동적 침투 |
| TS-REFL-003 | `_test.go` 내 `//go:linkname` | Review | 타 패키지 비공개 심볼 |

제외(위양성 방지): `reflect.DeepEqual`/`TypeOf`/단독 `ValueOf` — 정당.

구현 전제(rulebook §E-001 정정 사항):
- 표적이 `*_test.go`인데 현재 인덱서는 이를 제외 → **별도 test-file walk 필요**. 재사용 가능한 건
  `.tsmignore` 파서(E-002)뿐, 파일 적격 필터(E-001)는 정반대.
- Go는 go/ast로 정밀 탐지 가능. 다국어 확장(C#/Java/TS)은 라인 파서라 substring 위양성 노출(E-004).
- verdict 통합: tsma엔 `REVIEW` 상태 부재(S-001) → reins의 `LevelReview`/`State REVIEW`로 해결.

기준선(위양성 0 회귀): tsma 자체 테스트 = `unsafe` 0건, `//go:linkname` 0건,
`reflect` 8건 전부 `DeepEqual`. → TS-REFL-001/002/003 모두 현재 발화 0이어야 정상.

---

## 7. reins 포팅 시 규칙 재배치 요약

- **G-002**(커버리지 100%) = 중심 `gate.Rule`(LevelFail). 나머지 outcome은 reins 래칫이 자동 처리
  (PASS lock / MaxTries→DONE / FAIL 재시도).
- **G-001**(테스트 실패) = 상류 `LevelFail` 룰. graph 백엔드 쓰면 G-002를 `Supersedes`(테스트가
  깨지면 커버리지 판정 무의미).
- **TS-REFL-***  = 추가 `LevelReview` 룰. reins가 `rules` 카탈로그·`REVIEW` 상태를 공짜로 공급.
- **E-***/**M-***/**P-*** = 도메인 로직 → `Seed`/`Prepare`/`Render` 뒤로 이동(삭제 아님).
- **부작용**(테스트 실행·커버리지 측정) = `Prepare`에 격리. reins `ground`(HTTP/DNS 지향)는 부적합 →
  Prepare에서 도구 실행 후 `gate.Context`에 결과 주입, 룰은 순수 유지.
- **신규 능력**: `Options{Loop}`+`pkg/llm`(`claude:opus` 등)로 무인 구동(LLM 생성→게이트 판정→재시도)
  — tsma가 현재 미보유.

---

## 8. TANGEUL 매핑 — 레거시 → gate.Rule → gate.md 노드

§2·§6의 레거시 판정 로직은 이제 `internal/tsmagate/gate.md`(TANGEUL 판정 문서)로 **위상**이
선언된다. 술어 본문(`rule_*.go`)은 무변경이고, 문서는 15개 규칙의 위상(일반/반박)·레벨(실패/검토)·
우선순위(무효다/배제한다)만 코드 없이 드러낸다. 아래는 레거시 ID 기준 3층 대응이다
(이름 매핑 SSOT: `plans/tangeul/Phase002-gate-doc.md` §2).

| 레거시 | gate.Rule ID | gate.md 노드(한국어 이름) | 위상 |
|---|---|---|---|
| §2 G-001 | `tests-must-pass` | `테스트 실패` | 실패를 보고하는 반박 규칙, `커버리지 미달` 배제 |
| §2 G-002/G-004 | `branch-coverage-below-100` | `커버리지 미달` | 실패를 보고하는 반박 규칙 |
| §6 TS-REFL-001 | `TS-REFL-001` | `Go 언세이프` | 검토를 보고하는 반박 규칙 |
| §6 TS-REFL-002 | `TS-REFL-002` | `Go 리플렉트` | 〃 |
| §6 TS-REFL-003 | `TS-REFL-003` | `Go 링크네임` | 〃 |
| §6 TS-REFL-TS-001~003 | 동일 ID | `TS 애니 캐스트` / `TS 리플렉트` / `TS 오운 프로퍼티` | 〃 |
| §6 TS-REFL-JV-001~002 | 동일 ID | `자바 리플렉트` / `자바 셋액세서블` | 〃 |
| §6 TS-REFL-CS-001~002 | 동일 ID | `C샵 리플렉트` / `C샵 리플렉트 인포` | 〃 |
| §6 TS-REFL-RS-001~003 | 동일 ID | `러스트 언세이프` / `러스트 트랜스뮤트` / `러스트 포인터` | 〃 |

- **일반 규칙 `제출 유효`**: gate.md의 유일한 일반 규칙. 위 15개 반박 규칙이 전부 이 하나를
  `무효다`로 무효화한다 — 어느 하나라도 발화하면 제출은 유효하지 않다.
- **확정 배제 사슬(Phase003 패리티 확정본)**: `테스트 실패` → `커버리지 미달` **단일 간선 하나뿐**.
  이는 legacy `rule_branch_coverage.go`의 `!TestFailed` Go 가드를 감사 표면에 드러낸 것이다.
  smell 13개(TS-REFL-*)는 **배제하지 않는다** — legacy에서 smell은 TestFailed와 무관하게
  발화·기록되므로, 배제하면 Facts가 사라져 패리티가 깨진다. 따라서 Fail 2개는 이 한 간선,
  smell 13개는 병렬 미배제 Review이고, legacy 평탄 판정과 Outcome·RootCause·Facts가 완전 동등하다.
- **⚠️ 경고**: `RootCause`/`RuleSystem`/`EscalateOn` 키는 **gate.Rule ID**(`tests-must-pass`,
  `branch-coverage-below-100` 등)이며 **gate.md의 한국어 문서 이름이 아니다**. `RulePred`가
  `RuleMeta.ID`를 승계하므로 `loop_config.go`의 `EscalateOn: ["branch-coverage-below-100"]`
  등은 무변경이다. 문서 표시 이름(`커버리지 미달`)을 키로 쓰지 말 것.
