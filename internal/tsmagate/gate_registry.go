//ff:func feature=gate type=helper control=sequence
//ff:what gateRegistry: gate.md의 `검사`.<규칙> 심볼 15개를 rule_*.go의 미노출 술어 변수에 바인딩한 tangeul.Registry를 만든다. RulePred가 RuleMeta.ID를 승계하므로 RootCause·RuleSystem·EscalateOn 키(tests-must-pass 등)는 무변경이다. tsma 술어는 needs가 전혀 없어(전부 tier-0) RulePred에 needs 인자를 주지 않는다 — Provider/Resolver nil의 전제.

package tsmagate

import "github.com/park-jun-woo/reins/pkg/tangeul"

// gateRegistry binds gate.md's `검사`.<X> symbols to the 15 rule predicates.
// RulePred inherits RuleMeta.ID, so RootCause/RuleSystem/EscalateOn keys are
// unchanged. tsma's predicates declare no ground needs (all tier-0), so no
// needs args are passed — the premise for GateOptions' nil Provider/Resolver.
func gateRegistry() *tangeul.Registry {
	reg := tangeul.NewRegistry()
	_ = reg.RegisterPred("검사", "TestsMustPass", tangeul.RulePred(testsMustPass))
	_ = reg.RegisterPred("검사", "BranchCoverageBelow100", tangeul.RulePred(branchCoverageBelow100))
	_ = reg.RegisterPred("검사", "GoUnsafe", tangeul.RulePred(unsafeInTest))
	_ = reg.RegisterPred("검사", "GoReflectDynamic", tangeul.RulePred(reflectDynamicInTest))
	_ = reg.RegisterPred("검사", "GoLinkname", tangeul.RulePred(linknameInTest))
	_ = reg.RegisterPred("검사", "TsAsAny", tangeul.RulePred(tsAsAnyInTest))
	_ = reg.RegisterPred("검사", "TsReflect", tangeul.RulePred(tsReflectInTest))
	_ = reg.RegisterPred("검사", "TsOwnProperty", tangeul.RulePred(tsOwnPropertyInTest))
	_ = reg.RegisterPred("검사", "JavaReflect", tangeul.RulePred(javaReflectInTest))
	_ = reg.RegisterPred("검사", "JavaSetAccessible", tangeul.RulePred(javaSetAccessibleInTest))
	_ = reg.RegisterPred("검사", "CsReflect", tangeul.RulePred(csReflectInTest))
	_ = reg.RegisterPred("검사", "CsReflectInfo", tangeul.RulePred(csReflectInfoInTest))
	_ = reg.RegisterPred("검사", "RsUnsafe", tangeul.RulePred(rsUnsafeInTest))
	_ = reg.RegisterPred("검사", "RsTransmute", tangeul.RulePred(rsTransmuteInTest))
	_ = reg.RegisterPred("검사", "RsPtr", tangeul.RulePred(rsPtrInTest))
	return reg
}
