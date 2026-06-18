//ff:func feature=gate type=helper control=selection
//ff:what commentTokensFor: 누적 엔진(merge_canonical)이 함수별 블록 마커를 그 언어의 주석 문법으로 찍도록, lang→commentTokens(line 주석 접두사 + 헤더 판정 키워드)를 돌려주는 언어별 어댑터다. Go/TS/Java/C#/Rust는 `//`, Python은 `#`. header 판정 키워드(package/import/use/using/from·require)는 헤더 합집합 머지(splitHeader)가 본문과 import 영역을 가르는 데 쓴다. 미지원 언어는 Go 토큰으로 폴백(기존 generic write가 Go 폴백을 따르던 것과 동형).

package tsmagate

// commentTokensFor returns the comment/header facts for lang. Unknown languages
// fall back to the Go tokens, matching the prior generic-write Go fallback.
func commentTokensFor(lang string) commentTokens {
	switch lang {
	case "python":
		return commentTokens{
			line:           "#",
			headerPrefixes: []string{"import ", "from ", "import\t", "from\t"},
		}
	case "typescript":
		return commentTokens{
			line:           "//",
			headerPrefixes: []string{"import ", "import\t", "import{", "export {"},
		}
	case "java":
		return commentTokens{
			line:           "//",
			headerPrefixes: []string{"package ", "import ", "package\t", "import\t"},
		}
	case "csharp":
		return commentTokens{
			line:           "//",
			headerPrefixes: []string{"using ", "namespace ", "using\t", "namespace\t"},
		}
	case "rust":
		return commentTokens{
			line:           "//",
			headerPrefixes: []string{"use ", "extern ", "use\t", "extern\t"},
		}
	default: // go and unknown
		return commentTokens{
			line:           "//",
			headerPrefixes: []string{"package ", "import ", "package\t", "import\t"},
		}
	}
}
