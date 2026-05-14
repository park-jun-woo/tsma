package model

import "testing"

func TestFindFunctionByQualifiedName(t *testing.T) {
	s := &Session{
		Functions: []Function{
			{QualifiedName: "pkg.Login", Name: "Login", Status: StatusTodo},
			{QualifiedName: "pkg.Signup", Name: "Signup", Status: StatusDone},
		},
	}

	f := s.FindFunction("pkg.Login")
	if f == nil {
		t.Fatal("FindFunction(\"pkg.Login\") returned nil")
	}
	if f.Name != "Login" {
		t.Errorf("Name = %q, want %q", f.Name, "Login")
	}
}

func TestFindFunctionByName(t *testing.T) {
	s := &Session{
		Functions: []Function{
			{QualifiedName: "pkg.Login", Name: "Login", Status: StatusTodo},
			{QualifiedName: "pkg.Signup", Name: "Signup", Status: StatusDone},
		},
	}

	f := s.FindFunction("Signup")
	if f == nil {
		t.Fatal("FindFunction(\"Signup\") returned nil")
	}
	if f.QualifiedName != "pkg.Signup" {
		t.Errorf("QualifiedName = %q, want %q", f.QualifiedName, "pkg.Signup")
	}
}

func TestFindFunctionNotFound(t *testing.T) {
	s := &Session{
		Functions: []Function{
			{QualifiedName: "pkg.Login", Name: "Login", Status: StatusTodo},
		},
	}

	f := s.FindFunction("Nonexistent")
	if f != nil {
		t.Errorf("FindFunction(\"Nonexistent\") = %v, want nil", f)
	}
}

func TestFindFunctionEmptySession(t *testing.T) {
	s := &Session{}
	f := s.FindFunction("anything")
	if f != nil {
		t.Errorf("FindFunction on empty session = %v, want nil", f)
	}
}

func TestFindFunctionReturnsMutablePointer(t *testing.T) {
	s := &Session{
		Functions: []Function{
			{QualifiedName: "pkg.Login", Name: "Login", Status: StatusTodo},
		},
	}

	f := s.FindFunction("Login")
	if f == nil {
		t.Fatal("FindFunction returned nil")
	}

	f.Status = StatusPass
	if s.Functions[0].Status != StatusPass {
		t.Errorf("mutation via pointer did not propagate: got %q, want %q", s.Functions[0].Status, StatusPass)
	}
}

func TestFindFunctionQualifiedNamePriority(t *testing.T) {
	// When a function's Name matches another function's QualifiedName,
	// the first match wins (QualifiedName is checked first).
	s := &Session{
		Functions: []Function{
			{QualifiedName: "Signup", Name: "AliasSignup", Status: StatusTodo},
			{QualifiedName: "pkg.Signup", Name: "Signup", Status: StatusDone},
		},
	}

	f := s.FindFunction("Signup")
	if f == nil {
		t.Fatal("FindFunction(\"Signup\") returned nil")
	}
	// First function's QualifiedName == "Signup", so it should match first.
	if f.Name != "AliasSignup" {
		t.Errorf("Name = %q, want %q (first match by QualifiedName)", f.Name, "AliasSignup")
	}
}
