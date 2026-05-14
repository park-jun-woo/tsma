package model

import "testing"

func TestSummaryStructFields(t *testing.T) {
	s := Summary{
		Total: 10,
		Pass:  5,
		Done:  3,
		Todo:  2,
	}

	if s.Total != 10 {
		t.Errorf("Total = %d, want 10", s.Total)
	}
	if s.Pass != 5 {
		t.Errorf("Pass = %d, want 5", s.Pass)
	}
	if s.Done != 3 {
		t.Errorf("Done = %d, want 3", s.Done)
	}
	if s.Todo != 2 {
		t.Errorf("Todo = %d, want 2", s.Todo)
	}
}

func TestSummaryZeroValue(t *testing.T) {
	var s Summary
	if s.Total != 0 {
		t.Errorf("zero Total = %d, want 0", s.Total)
	}
	if s.Pass != 0 {
		t.Errorf("zero Pass = %d, want 0", s.Pass)
	}
	if s.Done != 0 {
		t.Errorf("zero Done = %d, want 0", s.Done)
	}
	if s.Todo != 0 {
		t.Errorf("zero Todo = %d, want 0", s.Todo)
	}
}
