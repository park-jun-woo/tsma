package coverage

import "testing"

func TestMatchesCoverageKey(t *testing.T) {
	tests := []struct {
		name          string
		key           string
		normalizedRel string
		relFile       string
		projectRoot   string
		want          bool
	}{
		{
			name:          "exact match",
			key:           "src/handler.ts",
			normalizedRel: "src/handler.ts",
			relFile:       "src/handler.ts",
			projectRoot:   "/project",
			want:          true,
		},
		{
			name:          "suffix match with separator",
			key:           "/absolute/path/src/handler.ts",
			normalizedRel: "src/handler.ts",
			relFile:       "src/handler.ts",
			projectRoot:   "/project",
			want:          true,
		},
		{
			name:          "suffix match without leading separator",
			key:           "path/src/handler.ts",
			normalizedRel: "src/handler.ts",
			relFile:       "src/handler.ts",
			projectRoot:   "/project",
			want:          true,
		},
		{
			name:          "absolute match via project root",
			key:           "/project/src/handler.ts",
			normalizedRel: "src/handler.ts",
			relFile:       "src/handler.ts",
			projectRoot:   "/project",
			want:          true,
		},
		{
			name:          "no match",
			key:           "src/other.ts",
			normalizedRel: "src/handler.ts",
			relFile:       "src/handler.ts",
			projectRoot:   "/project",
			want:          false,
		},
		{
			name:          "suffix match even with empty project root",
			key:           "/project/src/handler.ts",
			normalizedRel: "src/handler.ts",
			relFile:       "src/handler.ts",
			projectRoot:   "",
			want:          true,
		},
		{
			name:          "suffix direct match",
			key:           "src/handler.ts",
			normalizedRel: "src/handler.ts",
			relFile:       "src/handler.ts",
			projectRoot:   "",
			want:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchesCoverageKey(tt.key, tt.normalizedRel, tt.relFile, tt.projectRoot)
			if got != tt.want {
				t.Errorf("matchesCoverageKey(%q, %q, %q, %q) = %v, want %v",
					tt.key, tt.normalizedRel, tt.relFile, tt.projectRoot, got, tt.want)
			}
		})
	}
}
