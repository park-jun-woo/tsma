# juicer tsma 테스트 리포트

- 프로젝트: `~/.clari/repos/fullend/juicer/`
- 모듈: `github.com/park-jun-woo/juicer`
- 실행일: 2026-05-15
- tsma 결과: TODO 0 / 140

## 요약

| 항목 | 수치 |
|---|---|
| 함수 | 140 |
| PASS (100% 커버리지) | 114 (81.4%) |
| DONE (실용적 최대) | 26 (18.6%) |
| TODO | 0 |
| 커버리지 평균 | 96% |
| 테스트 파일 | 11개 |
| 패키지 | 4개 (main, ddl, scanner, sqls) |

## 패키지별

### ddl (16 함수)

| 함수 | 상태 | 커버리지 |
|---|---|---|
| Run | PASS | 100% |
| Parse | PASS | 100% |
| splitStatements | PASS | 100% |
| stripLeadingComments | PASS | 100% |
| applyStatement | PASS | 100% |
| stripInlineComments | PASS | 100% |
| applyCreateTable | DONE | 82% |
| extractParenBody | PASS | 100% |
| splitTopLevel | PASS | 100% |
| extractColumnName | PASS | 100% |
| applyAlterTable | DONE | 94% |
| splitAlterClauses | PASS | 100% |
| hasColumn | PASS | 100% |
| removeColumn | PASS | 100% |
| applyDropIndex | PASS | 100% |
| applyCreateIndex | PASS | 100% |
| Render | PASS | 100% |
| renderTable | PASS | 100% |
| cleanLine | PASS | 100% |
| WriteFiles | PASS | 100% |

PASS: 18, DONE: 2

### main (5 함수)

| 함수 | 상태 | 커버리지 |
|---|---|---|
| runScan | DONE | 79% |
| runDDL | DONE | 75% |
| runSQL | DONE | 43% |
| runSQLNext | DONE | 50% |
| printUsage | PASS | 100% |

PASS: 1, DONE: 4. CLI 진입점은 `os.Exit`, 플래그 파싱 등으로 커버리지 한계.

### scanner (60 함수)

#### handler.go

| 함수 | 상태 | 커버리지 |
|---|---|---|
| buildFuncIndex | PASS | 100% |
| analyzeHandlers | DONE | 88% |
| findInfoForExpr | PASS | 100% |
| analyzeExpr | DONE | 43% |
| lookupFunc | PASS | 100% |
| ginCtxParamName | PASS | 100% |
| isGinContextType | PASS | 100% |
| scanBody | PASS | 100% |
| checkOneDepthCall | DONE | 50% |
| resolveCallerArgs | DONE | 90% |
| isGinContextTypeInfo | PASS | 100% |
| isIntKind | PASS | 100% |
| handleBind | PASS | 100% |
| bindVarName | PASS | 100% |
| handleQuery | PASS | 100% |
| handlePathParam | PASS | 100% |
| handleForm | PASS | 100% |
| handleFile | PASS | 100% |
| handleResponse | PASS | 100% |
| resolveStatusCode | DONE | 79% |
| constToString | PASS | 100% |
| stringLitValue | PASS | 100% |
| exprString | DONE | 94% |
| ensureRequest | PASS | 100% |

#### openapi.go

| 함수 | 상태 | 커버리지 |
|---|---|---|
| ToOpenAPI | PASS | 100% |
| buildSpecNode | PASS | 100% |
| buildOperation | PASS | 100% |
| buildRequestBody | PASS | 100% |
| bodySchema | PASS | 100% |
| buildResponses | PASS | 100% |
| responseSchema | PASS | 100% |
| fieldsToSchema | PASS | 100% |
| fieldToProperty | PASS | 100% |
| goTypeToOpenAPI | PASS | 100% |
| goTypeFormat | PASS | 100% |
| isRequired | PASS | 100% |
| ginPathToOpenAPI | PASS | 100% |
| generateOperationID | PASS | 100% |
| pathMethodToOperationID | PASS | 100% |
| deduplicateEndpoints | PASS | 100% |
| richness | PASS | 100% |
| pickBestResponse | PASS | 100% |
| statusDescription | PASS | 100% |
| lcFirst | PASS | 100% |
| sortedPaths | PASS | 100% |
| toYAMLNode | DONE | 95% |

#### route.go

| 함수 | 상태 | 커버리지 |
|---|---|---|
| extractRoutes | PASS | 100% |
| scanFile | PASS | 100% |
| ginPkgName | PASS | 100% |
| registerParams | PASS | 100% |
| isGinRouterType | PASS | 100% |
| walkStmts | PASS | 100% |
| processAssign | PASS | 100% |
| tryRouteCall | PASS | 100% |
| tryUseCall | PASS | 100% |
| isGinInit | PASS | 100% |
| identName | PASS | 100% |
| exprName | PASS | 100% |
| joinPath | PASS | 100% |
| unquote | PASS | 100% |
| extractPathString | PASS | 100% |
| collectStringParts | PASS | 100% |
| pathParams | PASS | 100% |

#### scanner.go

| 함수 | 상태 | 커버리지 |
|---|---|---|
| Scan | DONE | 75% |

#### types.go

| 함수 | 상태 | 커버리지 |
|---|---|---|
| resolveBindType | PASS | 100% |
| resolveExprType | DONE | 94% |
| resolveType | PASS | 100% |
| extractFields | PASS | 100% |
| resolveEmbedded | PASS | 100% |
| resolveNestedFields | PASS | 100% |
| resolveResponseType | DONE | 75% |
| isGinH | PASS | 100% |
| extractGinHFields | PASS | 100% |
| inferValueType | PASS | 100% |
| formatType | PASS | 100% |
| unwrapPointer | PASS | 100% |

#### output.go

| 함수 | 상태 | 커버리지 |
|---|---|---|
| Render | PASS | 100% |

scanner 합계 — PASS: 48, DONE: 12

### sqls (39 함수)

#### next.go

| 함수 | 상태 | 커버리지 |
|---|---|---|
| RunNext | DONE | 81% |
| RunStatus | PASS | 100% |
| RunList | PASS | 100% |
| RunSkip | DONE | 89% |
| RunReset | PASS | 100% |
| createSession | DONE | 92% |
| firstTODO | PASS | 100% |
| countStatus | PASS | 100% |
| toQueryName | PASS | 100% |
| queryExists | PASS | 100% |
| runSqlcGenerate | DONE | 60% |
| printSkeleton | PASS | 100% |
| sqlcHint | PASS | 100% |
| formatSlice | PASS | 100% |

#### session.go

| 함수 | 상태 | 커버리지 |
|---|---|---|
| SessionDir | PASS | 100% |
| sessionPath | PASS | 100% |
| SessionExists | PASS | 100% |
| LoadSession | PASS | 100% |
| SaveSession | DONE | 86% |
| DeleteSession | PASS | 100% |

#### sqls.go

| 함수 | 상태 | 커버리지 |
|---|---|---|
| Extract | DONE | 93% |
| parseFile | DONE | 92% |
| receiverTypeName | PASS | 100% |
| detectCRUD | DONE | 93% |
| refineCRUD | PASS | 100% |
| refineCRUDFromAST | DONE | 95% |
| collectSQLFragments | PASS | 100% |
| collectInlineSQLArgs | PASS | 100% |
| normalizeWhitespace | PASS | 100% |
| extractTables | PASS | 100% |
| extractTablesFromSQL | PASS | 100% |
| appendUnique | PASS | 100% |
| detectDynamic | DONE | 94% |
| extractParams | PASS | 100% |
| extractReturns | PASS | 100% |
| typeString | PASS | 100% |
| RenderYAML | PASS | 100% |
| RenderJSON | PASS | 100% |

sqls 합계 — PASS: 28, DONE: 11

## DONE 사유 분류

| 사유 | 함수 수 | 예시 |
|---|---|---|
| `go/types` + `go/packages.Load` 의존 | 8 | analyzeExpr, checkOneDepthCall, resolveStatusCode, Scan |
| `os.Exit` / CLI 진입점 | 4 | runScan, runDDL, runSQL, runSQLNext |
| 외부 명령 (`sqlc`) 실행 | 2 | runSqlcGenerate, RunNext |
| 복합 AST 분기 도달 불가 | 7 | applyCreateTable, applyAlterTable, exprString, toYAMLNode |
| 파일 I/O 에러 경로 | 5 | SaveSession, Extract, parseFile, detectCRUD, detectDynamic |
