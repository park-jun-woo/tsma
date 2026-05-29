# testdata 픽스처 레이아웃

C#/Java/Rust 언어 지원 추가 작업(Phase 1~3)에서 사용할 테스트 픽스처의
디렉터리 구조 규칙을 정의한다. **실제 픽스처 파일은 각 Phase에서 추가한다.**
Phase 0 시점에는 이 레이아웃 규약만 확정한다.

## 분류 1 — 샘플 프로젝트 (E2E / detect / index / match 용)

각 언어별 최소 샘플 프로젝트(함수 2~3개 + 테스트)를 언어 디렉터리 아래 둔다.
`tsma next` 전체 루프(detect → index → match → run → coverage) 검증에 사용한다.

```
testdata/
  rust/
    sample-project/        # Cargo.toml 마커, fn/pub fn/impl, in-file #[cfg(test)] + tests/*.rs
  java/
    sample-maven/          # pom.xml 마커, src/main/java/... ↔ src/test/java/...
    sample-gradle/         # build.gradle(.kts) 마커 (여건 시)
  csharp/
    sample-project/        # *.csproj/*.sln 마커, Foo.cs ↔ FooTests.cs / *.Tests 프로젝트
```

규약:
- 각 샘플 프로젝트 루트에 해당 언어의 detect 마커 파일을 반드시 둔다.
- 소스/테스트 파일 경로는 indexer/matcher가 사용하는 **projectRoot 기준 상대경로**
  계약을 그대로 따른다.

## 분류 2 — 커버리지 도구 출력 샘플 (parser 단위 테스트 용)

실제 커버리지 도구가 뱉은 출력물을 그대로 커밋해 파서 단위 테스트의 입력으로 쓴다.
파서는 환경 비의존으로 검증 가능해야 하므로, 도구를 실행하지 않고 이 파일들을 읽어 파싱한다.

```
testdata/
  rust/
    coverage/
      llvm-cov.json        # cargo llvm-cov --json 출력 (parse_llvm_cov)
  java/
    coverage/
      jacoco.xml           # mvn test jacoco:report → target/site/jacoco/jacoco.xml (parse_jacoco)
  csharp/
    coverage/
      cobertura.xml        # dotnet test --collect:"XPlat Code Coverage" 출력 (parse_cobertura)
```

규약:
- 파일명은 위와 같이 고정한다. 케이스가 더 필요하면 접미사를 붙인다
  (예: `llvm-cov.empty.json`, `jacoco.branchless.xml`).
- 픽스처는 실제 도구 출력을 최소 가공으로 보존한다(라인 범위 매핑 검증을 위해
  메서드/라인 정보가 남아 있어야 한다).

## 참고

- 기존 Go/TS/Python 테스트는 `t.TempDir()`로 임시 파일을 생성하며 커밋된 픽스처를
  쓰지 않는다. 신규 언어는 detect/index/match는 동일하게 `t.TempDir()`를 쓸 수 있으나,
  커버리지 파서는 도구 출력 재현이 어렵고 형식 안정성 검증이 필요하므로 위
  `*/coverage/` 픽스처를 커밋해 사용한다.
- 디렉터리명 lang 키: `rust` / `java` / `csharp` (codebook.yaml `optional.lang`과 일치).
</content>
</invoke>
