package endpoint

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestPyDjangoDetectorFunctionView(t *testing.T) {
	dir := t.TempDir()

	urlsCode := `from django.urls import path
from . import views

urlpatterns = [
    path('users/', views.list_users),
]
`
	if err := os.WriteFile(filepath.Join(dir, "urls.py"), []byte(urlsCode), 0o644); err != nil {
		t.Fatal(err)
	}

	viewsCode := `from django.http import JsonResponse

def list_users(request):
    return JsonResponse([], safe=False)
`
	if err := os.WriteFile(filepath.Join(dir, "views.py"), []byte(viewsCode), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &PyDjangoDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 1 {
		t.Fatalf("got %d endpoints, want 1", len(endpoints))
	}

	ep := endpoints[0]
	if ep.Name != "list_users" {
		t.Errorf("name = %q, want list_users", ep.Name)
	}
	if ep.Method != "ANY" {
		t.Errorf("method = %q, want ANY (function-based view)", ep.Method)
	}
	if ep.Path != "users/" {
		t.Errorf("path = %q, want users/", ep.Path)
	}
	if ep.Status != "todo" {
		t.Errorf("status = %q, want todo", ep.Status)
	}
}

func TestPyDjangoDetectorClassBasedView(t *testing.T) {
	dir := t.TempDir()

	urlsCode := `from django.urls import path
from . import views

urlpatterns = [
    path('users/', views.UserView),
]
`
	if err := os.WriteFile(filepath.Join(dir, "urls.py"), []byte(urlsCode), 0o644); err != nil {
		t.Fatal(err)
	}

	viewsCode := `from django.views import View

class UserView(View):
    def get(self, request):
        return JsonResponse([])

    def post(self, request):
        return JsonResponse({})
`
	if err := os.WriteFile(filepath.Join(dir, "views.py"), []byte(viewsCode), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &PyDjangoDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 2 {
		t.Fatalf("got %d endpoints, want 2 (get + post methods)", len(endpoints))
	}

	sort.Slice(endpoints, func(i, j int) bool {
		return endpoints[i].Method < endpoints[j].Method
	})

	// GET
	if endpoints[0].Method != "GET" {
		t.Errorf("first endpoint method = %q, want GET", endpoints[0].Method)
	}
	if endpoints[0].Name != "UserView.get" {
		t.Errorf("first endpoint name = %q, want UserView.get", endpoints[0].Name)
	}
	if endpoints[0].Path != "users/" {
		t.Errorf("first endpoint path = %q, want users/", endpoints[0].Path)
	}

	// POST
	if endpoints[1].Method != "POST" {
		t.Errorf("second endpoint method = %q, want POST", endpoints[1].Method)
	}
	if endpoints[1].Name != "UserView.post" {
		t.Errorf("second endpoint name = %q, want UserView.post", endpoints[1].Name)
	}
}

func TestPyDjangoDetectorMultipleURLPatterns(t *testing.T) {
	dir := t.TempDir()

	urlsCode := `from django.urls import path
from . import views

urlpatterns = [
    path('users/', views.list_users),
    path('orders/', views.list_orders),
    path('items/', views.list_items),
]
`
	if err := os.WriteFile(filepath.Join(dir, "urls.py"), []byte(urlsCode), 0o644); err != nil {
		t.Fatal(err)
	}

	viewsCode := `from django.http import JsonResponse

def list_users(request):
    return JsonResponse([])

def list_orders(request):
    return JsonResponse([])

def list_items(request):
    return JsonResponse([])
`
	if err := os.WriteFile(filepath.Join(dir, "views.py"), []byte(viewsCode), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &PyDjangoDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 3 {
		t.Fatalf("got %d endpoints, want 3", len(endpoints))
	}

	names := make(map[string]bool)
	for _, ep := range endpoints {
		names[ep.Name] = true
	}
	for _, name := range []string{"list_users", "list_orders", "list_items"} {
		if !names[name] {
			t.Errorf("missing endpoint: %s", name)
		}
	}
}

func TestPyDjangoDetectorRePath(t *testing.T) {
	dir := t.TempDir()

	urlsCode := `from django.urls import re_path
from . import views

urlpatterns = [
    re_path(r'^users/(?P<pk>\d+)/$', views.user_detail),
]
`
	if err := os.WriteFile(filepath.Join(dir, "urls.py"), []byte(urlsCode), 0o644); err != nil {
		t.Fatal(err)
	}

	viewsCode := `from django.http import JsonResponse

def user_detail(request, pk):
    return JsonResponse({})
`
	if err := os.WriteFile(filepath.Join(dir, "views.py"), []byte(viewsCode), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &PyDjangoDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 1 {
		t.Fatalf("got %d endpoints, want 1", len(endpoints))
	}

	ep := endpoints[0]
	if ep.Name != "user_detail" {
		t.Errorf("name = %q, want user_detail", ep.Name)
	}
}

func TestPyDjangoDetectorSkipsPycache(t *testing.T) {
	dir := t.TempDir()

	// Valid urls.py at root.
	urlsCode := `from django.urls import path
from . import views
urlpatterns = [
    path('hello/', views.hello),
]
`
	if err := os.WriteFile(filepath.Join(dir, "urls.py"), []byte(urlsCode), 0o644); err != nil {
		t.Fatal(err)
	}

	viewsCode := `def hello(request):
    return "hi"
`
	if err := os.WriteFile(filepath.Join(dir, "views.py"), []byte(viewsCode), 0o644); err != nil {
		t.Fatal(err)
	}

	// urls.py inside __pycache__ should be skipped.
	cacheDir := filepath.Join(dir, "__pycache__")
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		t.Fatal(err)
	}
	cachedURLs := `from django.urls import path
urlpatterns = [
    path('cached/', views.cached),
]
`
	if err := os.WriteFile(filepath.Join(cacheDir, "urls.py"), []byte(cachedURLs), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &PyDjangoDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 1 {
		t.Errorf("got %d endpoints, want 1 (__pycache__ should be skipped)", len(endpoints))
	}
}

func TestPyDjangoDetectorEmptyDir(t *testing.T) {
	dir := t.TempDir()

	d := &PyDjangoDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 0 {
		t.Errorf("got %d endpoints, want 0", len(endpoints))
	}
}

func TestPyDjangoDetectorHandlerLocation(t *testing.T) {
	dir := t.TempDir()

	urlsCode := `from django.urls import path
from . import views
urlpatterns = [
    path('users/', views.list_users),
]
`
	if err := os.WriteFile(filepath.Join(dir, "urls.py"), []byte(urlsCode), 0o644); err != nil {
		t.Fatal(err)
	}

	viewsCode := `from django.http import JsonResponse

def list_users(request):
    users = []
    return JsonResponse(users, safe=False)
`
	if err := os.WriteFile(filepath.Join(dir, "views.py"), []byte(viewsCode), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &PyDjangoDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 1 {
		t.Fatalf("got %d endpoints, want 1", len(endpoints))
	}

	ep := endpoints[0]
	if ep.Handler.File == "" {
		t.Error("handler file should not be empty")
	}
	if ep.Handler.StartLine == 0 {
		t.Error("handler start line should not be 0")
	}
	if ep.Handler.EndLine < ep.Handler.StartLine {
		t.Errorf("handler end line (%d) < start line (%d)", ep.Handler.EndLine, ep.Handler.StartLine)
	}
}

func TestPyDjangoDetectorClassMethodLocation(t *testing.T) {
	dir := t.TempDir()

	urlsCode := `from django.urls import path
from . import views
urlpatterns = [
    path('items/', views.ItemView),
]
`
	if err := os.WriteFile(filepath.Join(dir, "urls.py"), []byte(urlsCode), 0o644); err != nil {
		t.Fatal(err)
	}

	viewsCode := `from django.views import View

class ItemView(View):
    def get(self, request):
        return JsonResponse([])

    def post(self, request):
        return JsonResponse({})

    def delete(self, request):
        return JsonResponse({})
`
	if err := os.WriteFile(filepath.Join(dir, "views.py"), []byte(viewsCode), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &PyDjangoDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 3 {
		t.Fatalf("got %d endpoints, want 3", len(endpoints))
	}

	for _, ep := range endpoints {
		if ep.Handler.StartLine == 0 {
			t.Errorf("endpoint %s: start line should not be 0", ep.Name)
		}
		if ep.Handler.EndLine < ep.Handler.StartLine {
			t.Errorf("endpoint %s: end line (%d) < start line (%d)", ep.Name, ep.Handler.EndLine, ep.Handler.StartLine)
		}
	}

	methods := make(map[string]bool)
	for _, ep := range endpoints {
		methods[ep.Method] = true
	}
	for _, m := range []string{"GET", "POST", "DELETE"} {
		if !methods[m] {
			t.Errorf("missing method: %s", m)
		}
	}
}

func TestPyDjangoDetectorUnresolvedView(t *testing.T) {
	dir := t.TempDir()

	// urls.py references a view that doesn't exist in any file.
	urlsCode := `from django.urls import path
from . import views
urlpatterns = [
    path('ghost/', views.ghost_handler),
]
`
	if err := os.WriteFile(filepath.Join(dir, "urls.py"), []byte(urlsCode), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &PyDjangoDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	// Even unresolved views should still produce an endpoint with minimal info.
	if len(endpoints) != 1 {
		t.Fatalf("got %d endpoints, want 1 (unresolved view still recorded)", len(endpoints))
	}

	ep := endpoints[0]
	if ep.Name != "ghost_handler" {
		t.Errorf("name = %q, want ghost_handler", ep.Name)
	}
	if ep.Method != "ANY" {
		t.Errorf("method = %q, want ANY", ep.Method)
	}
}

func TestPyDjangoDetectorViewsInSubdir(t *testing.T) {
	dir := t.TempDir()

	appDir := filepath.Join(dir, "myapp")
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		t.Fatal(err)
	}

	urlsCode := `from django.urls import path
from . import views
urlpatterns = [
    path('products/', views.list_products),
]
`
	if err := os.WriteFile(filepath.Join(appDir, "urls.py"), []byte(urlsCode), 0o644); err != nil {
		t.Fatal(err)
	}

	viewsCode := `from django.http import JsonResponse

def list_products(request):
    return JsonResponse([])
`
	if err := os.WriteFile(filepath.Join(appDir, "views.py"), []byte(viewsCode), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &PyDjangoDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 1 {
		t.Fatalf("got %d endpoints, want 1", len(endpoints))
	}

	if endpoints[0].Name != "list_products" {
		t.Errorf("name = %q, want list_products", endpoints[0].Name)
	}
}

func TestExtractViewName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"views.list_users", "list_users"},
		{"views.UserView", "UserView"},
		{"list_users", "list_users"},
		{"views.UserView.as_view", "UserView"},
		{"myapp.views.handler", "handler"},
	}

	for _, tt := range tests {
		got := extractViewName(tt.input)
		if got != tt.want {
			t.Errorf("extractViewName(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

// Ensure sort import is used.
var _ = sort.Slice
