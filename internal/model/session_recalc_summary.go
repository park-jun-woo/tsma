//ff:func feature=model type=implementation control=iteration
//ff:what Recomputes the summary counts from endpoint statuses
package model

// RecalcSummary recomputes the summary from endpoints.
func (s *Session) RecalcSummary() {
	s.Summary = Summary{Total: len(s.Endpoints)}
	for _, ep := range s.Endpoints {
		switch ep.Status {
		case StatusDone:
			s.Summary.Done++
		case StatusPartial:
			s.Summary.Partial++
		default:
			s.Summary.Todo++
		}
	}
}
