package cli

import (
	"math"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestComputeAverageCoverage_mixed(t *testing.T) {
	sess := &model.Session{
		Functions: []model.Function{
			{Name: "A", Status: model.StatusPass, CoveragePct: 100},
			{Name: "B", Status: model.StatusDone, CoveragePct: 80},
			{Name: "C", Status: model.StatusTodo, CoveragePct: 0},
		},
	}
	avg := computeAverageCoverage(sess)
	if math.Abs(avg-90) > 0.01 {
		t.Errorf("expected 90, got %f", avg)
	}
}

func TestComputeAverageCoverage_allTodo(t *testing.T) {
	sess := &model.Session{
		Functions: []model.Function{
			{Name: "A", Status: model.StatusTodo},
			{Name: "B", Status: model.StatusTodo},
		},
	}
	avg := computeAverageCoverage(sess)
	if avg != 0 {
		t.Errorf("expected 0, got %f", avg)
	}
}

func TestComputeAverageCoverage_empty(t *testing.T) {
	sess := &model.Session{
		Functions: []model.Function{},
	}
	avg := computeAverageCoverage(sess)
	if avg != 0 {
		t.Errorf("expected 0, got %f", avg)
	}
}

func TestComputeAverageCoverage_allPass(t *testing.T) {
	sess := &model.Session{
		Functions: []model.Function{
			{Name: "A", Status: model.StatusPass, CoveragePct: 100},
			{Name: "B", Status: model.StatusPass, CoveragePct: 60},
			{Name: "C", Status: model.StatusPass, CoveragePct: 80},
		},
	}
	avg := computeAverageCoverage(sess)
	if math.Abs(avg-80) > 0.01 {
		t.Errorf("expected 80, got %f", avg)
	}
}

func TestComputeAverageCoverage_singleDone(t *testing.T) {
	sess := &model.Session{
		Functions: []model.Function{
			{Name: "A", Status: model.StatusDone, CoveragePct: 75},
		},
	}
	avg := computeAverageCoverage(sess)
	if math.Abs(avg-75) > 0.01 {
		t.Errorf("expected 75, got %f", avg)
	}
}
