//ff:func feature=model type=implementation control=iteration dimension=1
//ff:what Returns a pointer to the function matching the name or qualified name
package model

func (s *Session) FindFunction(name string) *Function {
	for i := range s.Functions {
		if s.Functions[i].QualifiedName == name || s.Functions[i].Name == name {
			return &s.Functions[i]
		}
	}
	return nil
}
