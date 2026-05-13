package model

import "testing"

func TestRecalcSummary(t *testing.T) {
	s := &Session{
		Endpoints: []Endpoint{
			{Name: "A", Status: StatusDone},
			{Name: "B", Status: StatusPartial},
			{Name: "C", Status: StatusTodo},
			{Name: "D", Status: StatusTodo},
		},
	}
	s.RecalcSummary()

	if s.Summary.Total != 4 {
		t.Errorf("Total = %d, want 4", s.Summary.Total)
	}
	if s.Summary.Done != 1 {
		t.Errorf("Done = %d, want 1", s.Summary.Done)
	}
	if s.Summary.Partial != 1 {
		t.Errorf("Partial = %d, want 1", s.Summary.Partial)
	}
	if s.Summary.Todo != 2 {
		t.Errorf("Todo = %d, want 2", s.Summary.Todo)
	}
}

func TestFindEndpoint(t *testing.T) {
	s := &Session{
		Endpoints: []Endpoint{
			{Name: "Login"},
			{Name: "Signup"},
		},
	}

	ep := s.FindEndpoint("Signup")
	if ep == nil {
		t.Fatal("FindEndpoint returned nil for existing endpoint")
	}
	if ep.Name != "Signup" {
		t.Errorf("Name = %q, want Signup", ep.Name)
	}

	ep = s.FindEndpoint("NotExist")
	if ep != nil {
		t.Error("FindEndpoint should return nil for non-existing endpoint")
	}
}
