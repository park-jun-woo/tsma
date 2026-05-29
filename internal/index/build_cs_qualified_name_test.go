package index

import "testing"

func TestBuildCsQualifiedNameWithBlockNamespace(t *testing.T) {
	scopes := []csScope{
		{typeName: "Com.Example"},
		{typeName: "Calculator"},
	}
	got := buildCsQualifiedName("", scopes, "Add")
	if got != "Com.Example.Calculator.Add" {
		t.Errorf("got %q, want Com.Example.Calculator.Add", got)
	}
}

func TestBuildCsQualifiedNameWithFileNamespace(t *testing.T) {
	scopes := []csScope{{typeName: "Calculator"}}
	got := buildCsQualifiedName("Com.Example", scopes, "Add")
	if got != "Com.Example.Calculator.Add" {
		t.Errorf("got %q, want Com.Example.Calculator.Add", got)
	}
}

func TestBuildCsQualifiedNameBare(t *testing.T) {
	got := buildCsQualifiedName("", nil, "Main")
	if got != "Main" {
		t.Errorf("got %q, want Main", got)
	}
}
