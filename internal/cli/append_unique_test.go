//ff:test feature=cli
package cli

import (
	"reflect"
	"testing"
)

func TestAppendUnique(t *testing.T) {
	tests := []struct {
		name  string
		dst   []string
		items []string
		want  []string
	}{
		{"empty both", nil, nil, nil},
		{"into empty", nil, []string{"a", "b"}, []string{"a", "b"}},
		{"skip existing in dst", []string{"a"}, []string{"a", "b"}, []string{"a", "b"}},
		{"skip duplicate within items", nil, []string{"a", "a", "b"}, []string{"a", "b"}},
		{"all already present", []string{"a", "b"}, []string{"a", "b"}, []string{"a", "b"}},
		{"preserve order", []string{"x"}, []string{"z", "y", "x", "w"}, []string{"x", "z", "y", "w"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := appendUnique(tt.dst, tt.items)
			if len(got) == 0 && len(tt.want) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendUnique(%v, %v) = %v, want %v", tt.dst, tt.items, got, tt.want)
			}
		})
	}
}
