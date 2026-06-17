//ff:type feature=gate type=model
//ff:what Definition은 tsma의 reins gate.Definition 구현체다. Seed/Render/Prepare/Rules 4메서드(각 파일)로 reins 래칫에 도메인 로직(인덱싱·매칭·테스트실행·커버리지측정·100%브랜치 게이트)을 끼운다. 가변 상태가 없어 단일 값을 모든 서브커맨드가 공유한다. New(new.go)가 생성하고, 컴파일 타임 assertion으로 reins 계약을 보장한다.

package tsmagate

import "github.com/park-jun-woo/reins/pkg/gate"

// Definition implements gate.Definition for tsma. It carries no mutable state;
// per-quest data rides in each Item's payload, so a single value is shared
// across all subcommands.
type Definition struct{}

// Ensure Definition satisfies the reins contract at compile time.
var _ gate.Definition = (*Definition)(nil)
