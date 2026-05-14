package index

import "testing"

func TestBuildQualifiedName(t *testing.T) {
	tests := []struct {
		name     string
		pkgDir   string
		receiver string
		funcName string
		want     string
	}{
		{
			name:     "package + receiver + function",
			pkgDir:   "internal/api",
			receiver: "Handler",
			funcName: "Login",
			want:     "internal/api.Handler.Login",
		},
		{
			name:     "package + function only",
			pkgDir:   "internal/api",
			receiver: "",
			funcName: "helperFunc",
			want:     "internal/api.helperFunc",
		},
		{
			name:     "root package with receiver",
			pkgDir:   "",
			receiver: "Server",
			funcName: "Start",
			want:     "Server.Start",
		},
		{
			name:     "root package function only",
			pkgDir:   "",
			receiver: "",
			funcName: "main",
			want:     "main",
		},
		{
			name:     "deep nested package",
			pkgDir:   "pkg/auth/service",
			receiver: "AuthSvc",
			funcName: "Validate",
			want:     "pkg/auth/service.AuthSvc.Validate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildQualifiedName(tt.pkgDir, tt.receiver, tt.funcName)
			if got != tt.want {
				t.Errorf("buildQualifiedName(%q, %q, %q) = %q, want %q",
					tt.pkgDir, tt.receiver, tt.funcName, got, tt.want)
			}
		})
	}
}
