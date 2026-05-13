//ff:func feature=model type=implementation control=iteration dimension=1
//ff:what Returns all functions matching a short name for ambiguity detection
package model

func (s *Session) FindAmbiguous(name string) []*Function {
	var matches []*Function
	for i := range s.Functions {
		if s.Functions[i].Name == name {
			matches = append(matches, &s.Functions[i])
		}
	}
	return matches
}
