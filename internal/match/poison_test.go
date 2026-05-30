package match

import "testing"

func TestPoison(t *testing.T) {
	types := map[string]string{"x": "GoFile"}
	poisoned := make(map[string]struct{})
	poison(types, poisoned, "x")
	if _, ok := types["x"]; ok {
		t.Errorf("poison must delete the tracked type for x")
	}
	if _, ok := poisoned["x"]; !ok {
		t.Errorf("poison must record x as poisoned")
	}
}
