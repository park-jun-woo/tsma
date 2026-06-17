//ff:func feature=gate type=helper control=sequence level=error lang=go
//ff:what buildLoopTestMatch: 생성 소스를 .tsma/test/gen-<item>.go(backing, gitignore)에 쓰고, 패키지 디렉터리 안의 충돌없는 가상경로(zzz_tsma_gen_test.go)→backing 절대경로 overlay를 구성해 TestMatch를 직접 만들어 돌려준다(MatchFunc 디스크 재매칭 우회). 소스 트리에는 한 글자도 안 쓴다 — 측정은 go test -overlay가 가상파일을 패키지에 끼워 수행한다. backing의 root-상대 경로도 함께 돌려(smell 스캔·materialize용). 쓰기 실패는 error로 전파해 Prepare가 TestFailed로 드러낸다.
package tsmagate

import (
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/park-jun-woo/tsma/internal/match"
)

// buildLoopTestMatch writes the generated source to a backing file under
// .tsma/test (gitignored) and constructs the overlay TestMatch by hand: the
// virtual path is a fresh collision-free _test.go inside the source package
// directory, mapped (absolute → absolute) to that backing file. The source tree
// is never written — go test -overlay splices the virtual file into the package
// at build time. It returns the match plus the backing's root-relative path (for
// the smell scan and materialize). MatchFunc's disk re-match is bypassed so a
// brand-new test still attributes. A write failure propagates as an error.
func buildLoopTestMatch(p funcPayload, it *quest.Item, src string, funcs []string) (match.TestMatch, string, error) {
	slug := strings.NewReplacer("/", "_", ".", "_", "*", "_", "(", "_", ")", "_", " ", "_", "[", "_", "]", "_").Replace(it.Key)
	backingRel := filepath.Join(".tsma", "test", "gen-"+slug+".go")
	if err := writeTestFile(p.Root, backingRel, src); err != nil {
		return match.TestMatch{}, "", err
	}
	virtualRel := filepath.Join(filepath.Dir(p.Fn.File), "zzz_tsma_gen_test.go")
	overlay := map[string]string{
		filepath.Join(p.Root, virtualRel): filepath.Join(p.Root, backingRel),
	}
	tm := match.TestMatch{Files: []string{virtualRel}, TestFuncs: funcs, Overlay: overlay}
	return tm, backingRel, nil
}
