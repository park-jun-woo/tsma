//ff:func feature=gate type=helper control=sequence
//ff:what New: tsma gate Definition 생성자. 가변 상태가 없으므로 빈 &Definition{}를 반환한다.

package tsmagate

// New returns a tsma gate Definition.
func New() *Definition { return &Definition{} }
