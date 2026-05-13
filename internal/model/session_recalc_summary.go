//ff:func feature=model type=implementation control=iteration dimension=1
//ff:what Recomputes the summary counts from function statuses
package model

func (s *Session) RecalcSummary() {
	s.Summary = Summary{}
	for _, fn := range s.Functions {
		if fn.Dead {
			continue
		}
		s.Summary.Testable++
		switch fn.Status {
		case StatusDone:
			s.Summary.Done++
		case StatusPartial:
			s.Summary.Partial++
		default:
			s.Summary.Todo++
		}
	}
}
