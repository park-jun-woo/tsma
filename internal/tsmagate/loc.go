//ff:func feature=gate type=helper control=sequence
//ff:what loc: Finding의 위치를 Fact.Where용 "file:line" 문자열로 렌더한다.

package tsmagate

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/smell"
)

// loc renders a Finding's location as "file:line" for Fact.Where.
func loc(f *smell.Finding) string {
	return fmt.Sprintf("%s:%d", f.File, f.Line)
}
