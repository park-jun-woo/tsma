package model

import (
	"testing"
	"time"
)

func TestSessionStructFields(t *testing.T) {
	now := time.Now()
	s := Session{
		Project:      "/tmp/myproject",
		Lang:         "go",
		CheckedAt:    now,
		CurrentIndex: 3,
		Functions: []Function{
			{Name: "Login", Status: StatusPass},
			{Name: "Signup", Status: StatusTodo},
		},
		Summary: Summary{Total: 2, Pass: 1, Todo: 1},
	}

	if s.Project != "/tmp/myproject" {
		t.Errorf("Project = %q, want %q", s.Project, "/tmp/myproject")
	}
	if s.Lang != "go" {
		t.Errorf("Lang = %q, want %q", s.Lang, "go")
	}
	if !s.CheckedAt.Equal(now) {
		t.Errorf("CheckedAt = %v, want %v", s.CheckedAt, now)
	}
	if s.CurrentIndex != 3 {
		t.Errorf("CurrentIndex = %d, want 3", s.CurrentIndex)
	}
	if len(s.Functions) != 2 {
		t.Errorf("Functions count = %d, want 2", len(s.Functions))
	}
	if s.Summary.Total != 2 {
		t.Errorf("Summary.Total = %d, want 2", s.Summary.Total)
	}
}

func TestSessionZeroValue(t *testing.T) {
	var s Session
	if s.Project != "" {
		t.Errorf("zero Project = %q, want empty", s.Project)
	}
	if s.Lang != "" {
		t.Errorf("zero Lang = %q, want empty", s.Lang)
	}
	if s.CurrentIndex != 0 {
		t.Errorf("zero CurrentIndex = %d, want 0", s.CurrentIndex)
	}
	if len(s.Functions) != 0 {
		t.Errorf("zero Functions count = %d, want 0", len(s.Functions))
	}
	if s.Summary.Total != 0 {
		t.Errorf("zero Summary.Total = %d, want 0", s.Summary.Total)
	}
}
