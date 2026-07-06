//ff:func feature=cli type=command control=sequence level=error
//ff:what tsma 엔트리포인트. reins NewQuestCmd로 표준 quest CLI(scan/next/submit/status/export/rules)를 조립하고, 도메인 게이트는 tsmagate.Definition만 끼운다. Options.Loop에 tsmagate.LoopOptions()를 실어 무인 generate→gate→retry `loop` 명령을 부착한다.

// Command tsma drives LLMs toward 100% branch coverage across languages, on top
// of the reins deterministic quest gate. The domain logic (seed/render/measure/
// rules) lives in internal/tsmagate; reins supplies the ratchet, command
// skeleton, aggregation, and export.
package main

import (
	"github.com/park-jun-woo/reins/pkg/cli"
	"github.com/park-jun-woo/tsma/internal/tsmagate"
)

func main() {
	def := tsmagate.New()
	cli.NewQuestCmd("tsma", def, cli.Options{
		Version: "0.5.1",
		Loop:    tsmagate.LoopOptions(),
		Gate:    tsmagate.GateOptions(),
	}).Execute()
}
