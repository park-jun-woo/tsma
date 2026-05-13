//ff:func feature=model type=implementation control=iteration dimension=1
//ff:what Recomputes the summary counts from function statuses
package model

func (s *Session) RecalcSummary() {
	s.Summary = Summary{Total: len(s.Functions)}
	for _, fn := range s.Functions {
		switch fn.Status {
		case StatusPass:
			s.Summary.Pass++
		case StatusDone:
			s.Summary.Done++
		default:
			s.Summary.Todo++
		}
	}
}
