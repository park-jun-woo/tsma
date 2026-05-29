package index

import "testing"

func TestUpdateLastNonEmpty(t *testing.T) {
	t.Run("updates when non-empty", func(t *testing.T) {
		last := 3
		updateLastNonEmpty(true, 10, &last)
		if last != 10 {
			t.Errorf("expected last=10, got %d", last)
		}
	})

	t.Run("leaves unchanged when empty", func(t *testing.T) {
		last := 3
		updateLastNonEmpty(false, 10, &last)
		if last != 3 {
			t.Errorf("expected last to stay 3, got %d", last)
		}
	})
}
