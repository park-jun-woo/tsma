//ff:func feature=model type=implementation control=iteration
//ff:what Returns a pointer to the named endpoint or nil if not found
package model

// FindEndpoint returns a pointer to the named endpoint, or nil.
func (s *Session) FindEndpoint(name string) *Endpoint {
	for i := range s.Endpoints {
		if s.Endpoints[i].Name == name {
			return &s.Endpoints[i]
		}
	}
	return nil
}
