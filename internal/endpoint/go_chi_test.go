package endpoint

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestGoChiDetector(t *testing.T) {
	dir := t.TempDir()

	code := `package main

import "github.com/go-chi/chi/v5"

func setupRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/users", ListUsers)
	r.Post("/users", CreateUser)
	r.Delete("/users/{id}", DeleteUser)
	return r
}

func ListUsers(w http.ResponseWriter, r *http.Request) {}
func CreateUser(w http.ResponseWriter, r *http.Request) {}
func DeleteUser(w http.ResponseWriter, r *http.Request) {}
`
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &GoChiDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 3 {
		t.Fatalf("got %d endpoints, want 3", len(endpoints))
	}

	expected := map[string]struct {
		method string
		path   string
	}{
		"ListUsers":  {method: "GET", path: "/users"},
		"CreateUser": {method: "POST", path: "/users"},
		"DeleteUser": {method: "DELETE", path: "/users/{id}"},
	}

	for _, ep := range endpoints {
		want, ok := expected[ep.Name]
		if !ok {
			t.Errorf("unexpected endpoint: %s", ep.Name)
			continue
		}
		if ep.Method != want.method {
			t.Errorf("endpoint %s: method = %q, want %q", ep.Name, ep.Method, want.method)
		}
		if ep.Path != want.path {
			t.Errorf("endpoint %s: path = %q, want %q", ep.Name, ep.Path, want.path)
		}
		if ep.Status != "todo" {
			t.Errorf("endpoint %s: status = %q, want todo", ep.Name, ep.Status)
		}
	}
}

func TestGoChiDetectorMethodNormalization(t *testing.T) {
	dir := t.TempDir()

	// Chi uses lowercase-initial methods (Get, Post, etc.); the detector should
	// normalize them to uppercase HTTP methods (GET, POST, etc.).
	code := `package main

import "github.com/go-chi/chi/v5"

func setup(r chi.Router) {
	r.Get("/a", HandlerA)
	r.Post("/b", HandlerB)
	r.Put("/c", HandlerC)
	r.Patch("/d", HandlerD)
	r.Delete("/e", HandlerE)
	r.Head("/f", HandlerF)
	r.Options("/g", HandlerG)
}

func HandlerA(w http.ResponseWriter, r *http.Request) {}
func HandlerB(w http.ResponseWriter, r *http.Request) {}
func HandlerC(w http.ResponseWriter, r *http.Request) {}
func HandlerD(w http.ResponseWriter, r *http.Request) {}
func HandlerE(w http.ResponseWriter, r *http.Request) {}
func HandlerF(w http.ResponseWriter, r *http.Request) {}
func HandlerG(w http.ResponseWriter, r *http.Request) {}
`
	if err := os.WriteFile(filepath.Join(dir, "routes.go"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &GoChiDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 7 {
		t.Fatalf("got %d endpoints, want 7", len(endpoints))
	}

	wantMethods := map[string]string{
		"HandlerA": "GET",
		"HandlerB": "POST",
		"HandlerC": "PUT",
		"HandlerD": "PATCH",
		"HandlerE": "DELETE",
		"HandlerF": "HEAD",
		"HandlerG": "OPTIONS",
	}

	for _, ep := range endpoints {
		wantMethod, ok := wantMethods[ep.Name]
		if !ok {
			t.Errorf("unexpected endpoint: %s", ep.Name)
			continue
		}
		if ep.Method != wantMethod {
			t.Errorf("endpoint %s: method = %q, want %q (normalized uppercase)", ep.Name, ep.Method, wantMethod)
		}
	}
}

func TestGoChiDetectorSkipsTestFiles(t *testing.T) {
	dir := t.TempDir()

	code := `package main
import "github.com/go-chi/chi/v5"
func setup(r chi.Router) {
	r.Get("/hello", Hello)
}
func Hello(w http.ResponseWriter, r *http.Request) {}
`
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	testCode := `package main
import "github.com/go-chi/chi/v5"
func setupTest(r chi.Router) {
	r.Get("/test", TestHello)
}
func TestHello(w http.ResponseWriter, r *http.Request) {}
`
	if err := os.WriteFile(filepath.Join(dir, "main_test.go"), []byte(testCode), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &GoChiDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 1 {
		t.Errorf("got %d endpoints, want 1 (test file should be skipped)", len(endpoints))
	}
}

func TestGoChiDetectorHandleFunc(t *testing.T) {
	dir := t.TempDir()

	code := `package main

import "github.com/go-chi/chi/v5"

func setup(r chi.Router) {
	r.HandleFunc("/health", HealthCheck)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {}
`
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &GoChiDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 1 {
		t.Fatalf("got %d endpoints, want 1", len(endpoints))
	}

	ep := endpoints[0]
	if ep.Name != "HealthCheck" {
		t.Errorf("name = %q, want HealthCheck", ep.Name)
	}
	// HandleFunc maps to the method name "HandleFunc" since there's no specific HTTP method.
	if ep.Method != "HandleFunc" {
		t.Errorf("method = %q, want HandleFunc", ep.Method)
	}
}

func TestGoChiDetectorEmptyDir(t *testing.T) {
	dir := t.TempDir()

	d := &GoChiDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 0 {
		t.Errorf("got %d endpoints, want 0", len(endpoints))
	}
}

func TestGoChiDetectorHandlerLocation(t *testing.T) {
	dir := t.TempDir()

	code := `package main

import "github.com/go-chi/chi/v5"

func setup(r chi.Router) {
	r.Get("/users", ListUsers)
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	// handler body
}
`
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &GoChiDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 1 {
		t.Fatalf("got %d endpoints, want 1", len(endpoints))
	}

	ep := endpoints[0]
	if ep.Handler.File != "main.go" {
		t.Errorf("handler file = %q, want main.go", ep.Handler.File)
	}
	if ep.Handler.StartLine == 0 {
		t.Error("handler start line should not be 0")
	}
	if ep.Handler.EndLine < ep.Handler.StartLine {
		t.Errorf("handler end line (%d) < start line (%d)", ep.Handler.EndLine, ep.Handler.StartLine)
	}
}

func TestGoChiDetectorSubdirectory(t *testing.T) {
	dir := t.TempDir()

	subDir := filepath.Join(dir, "handlers")
	if err := os.MkdirAll(subDir, 0o755); err != nil {
		t.Fatal(err)
	}

	code := `package handlers

import "github.com/go-chi/chi/v5"

func Setup(r chi.Router) {
	r.Get("/items", ListItems)
	r.Post("/items", CreateItem)
}

func ListItems(w http.ResponseWriter, r *http.Request) {}
func CreateItem(w http.ResponseWriter, r *http.Request) {}
`
	if err := os.WriteFile(filepath.Join(subDir, "items.go"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &GoChiDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 2 {
		t.Fatalf("got %d endpoints, want 2", len(endpoints))
	}

	sort.Slice(endpoints, func(i, j int) bool {
		return endpoints[i].Name < endpoints[j].Name
	})

	if endpoints[0].Name != "CreateItem" {
		t.Errorf("first endpoint name = %q, want CreateItem", endpoints[0].Name)
	}
	if endpoints[1].Name != "ListItems" {
		t.Errorf("second endpoint name = %q, want ListItems", endpoints[1].Name)
	}
}
