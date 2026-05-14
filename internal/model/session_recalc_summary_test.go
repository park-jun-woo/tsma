package model

import "testing"

func TestRecalcSummaryMixed(t *testing.T) {
	s := &Session{
		Functions: []Function{
			{Status: StatusPass},
			{Status: StatusDone},
			{Status: StatusTodo},
			{Status: StatusPass},
			{Status: ""},
		},
	}

	s.RecalcSummary()

	if s.Summary.Total != 5 {
		t.Errorf("Total = %d, want 5", s.Summary.Total)
	}
	if s.Summary.Pass != 2 {
		t.Errorf("Pass = %d, want 2", s.Summary.Pass)
	}
	if s.Summary.Done != 1 {
		t.Errorf("Done = %d, want 1", s.Summary.Done)
	}
	if s.Summary.Todo != 2 {
		t.Errorf("Todo = %d, want 2 (includes empty status)", s.Summary.Todo)
	}
}

func TestRecalcSummaryAllPass(t *testing.T) {
	s := &Session{
		Functions: []Function{
			{Status: StatusPass},
			{Status: StatusPass},
			{Status: StatusPass},
		},
	}

	s.RecalcSummary()

	if s.Summary.Total != 3 {
		t.Errorf("Total = %d, want 3", s.Summary.Total)
	}
	if s.Summary.Pass != 3 {
		t.Errorf("Pass = %d, want 3", s.Summary.Pass)
	}
	if s.Summary.Done != 0 {
		t.Errorf("Done = %d, want 0", s.Summary.Done)
	}
	if s.Summary.Todo != 0 {
		t.Errorf("Todo = %d, want 0", s.Summary.Todo)
	}
}

func TestRecalcSummaryAllTodo(t *testing.T) {
	s := &Session{
		Functions: []Function{
			{Status: StatusTodo},
			{Status: StatusTodo},
		},
	}

	s.RecalcSummary()

	if s.Summary.Total != 2 {
		t.Errorf("Total = %d, want 2", s.Summary.Total)
	}
	if s.Summary.Pass != 0 {
		t.Errorf("Pass = %d, want 0", s.Summary.Pass)
	}
	if s.Summary.Todo != 2 {
		t.Errorf("Todo = %d, want 2", s.Summary.Todo)
	}
}

func TestRecalcSummaryEmpty(t *testing.T) {
	s := &Session{}
	s.RecalcSummary()

	if s.Summary.Total != 0 {
		t.Errorf("Total = %d, want 0", s.Summary.Total)
	}
	if s.Summary.Pass != 0 {
		t.Errorf("Pass = %d, want 0", s.Summary.Pass)
	}
	if s.Summary.Done != 0 {
		t.Errorf("Done = %d, want 0", s.Summary.Done)
	}
	if s.Summary.Todo != 0 {
		t.Errorf("Todo = %d, want 0", s.Summary.Todo)
	}
}

func TestRecalcSummaryResetsOldValues(t *testing.T) {
	s := &Session{
		Summary: Summary{Total: 99, Pass: 50, Done: 30, Todo: 19},
		Functions: []Function{
			{Status: StatusPass},
		},
	}

	s.RecalcSummary()

	if s.Summary.Total != 1 {
		t.Errorf("Total = %d, want 1 (should be recalculated)", s.Summary.Total)
	}
	if s.Summary.Pass != 1 {
		t.Errorf("Pass = %d, want 1", s.Summary.Pass)
	}
	if s.Summary.Done != 0 {
		t.Errorf("Done = %d, want 0", s.Summary.Done)
	}
	if s.Summary.Todo != 0 {
		t.Errorf("Todo = %d, want 0", s.Summary.Todo)
	}
}

func TestRecalcSummaryUnknownStatusCountsAsTodo(t *testing.T) {
	s := &Session{
		Functions: []Function{
			{Status: "unknown_status"},
			{Status: "fail"},
			{Status: ""},
		},
	}

	s.RecalcSummary()

	if s.Summary.Todo != 3 {
		t.Errorf("Todo = %d, want 3 (all unknown statuses should be todo)", s.Summary.Todo)
	}
}
