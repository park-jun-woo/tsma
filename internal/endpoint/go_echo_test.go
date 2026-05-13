package endpoint

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestGoEchoDetector(t *testing.T) {
	dir := t.TempDir()

	code := `package main

import "github.com/labstack/echo/v4"

func setupRouter() *echo.Echo {
	e := echo.New()
	e.GET("/users", ListUsers)
	e.POST("/users", CreateUser)
	e.Add("PATCH", "/users/:id", UpdateUser)
	return e
}

func ListUsers(c echo.Context) error { return nil }
func CreateUser(c echo.Context) error { return nil }
func UpdateUser(c echo.Context) error { return nil }
`
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &GoEchoDetector{}
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
		"UpdateUser": {method: "PATCH", path: "/users/:id"},
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
		if ep.Handler.File == "" {
			t.Errorf("endpoint %s: handler file is empty", ep.Name)
		}
	}
}

func TestGoEchoDetectorSkipsTestFiles(t *testing.T) {
	dir := t.TempDir()

	code := `package main
import "github.com/labstack/echo/v4"
func setup(e *echo.Echo) {
	e.GET("/hello", Hello)
}
func Hello(c echo.Context) error { return nil }
`
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	testCode := `package main
import "github.com/labstack/echo/v4"
func setupTest(e *echo.Echo) {
	e.GET("/test", TestHandler)
}
func TestHandler(c echo.Context) error { return nil }
`
	if err := os.WriteFile(filepath.Join(dir, "main_test.go"), []byte(testCode), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &GoEchoDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 1 {
		t.Errorf("got %d endpoints, want 1 (test file should be skipped)", len(endpoints))
	}
}

func TestGoEchoDetectorSkipsVendor(t *testing.T) {
	dir := t.TempDir()

	vendorDir := filepath.Join(dir, "vendor")
	if err := os.MkdirAll(vendorDir, 0o755); err != nil {
		t.Fatal(err)
	}

	vendorCode := `package vendor
import "github.com/labstack/echo/v4"
func setup(e *echo.Echo) {
	e.GET("/vendored", VendoredHandler)
}
func VendoredHandler(c echo.Context) error { return nil }
`
	if err := os.WriteFile(filepath.Join(vendorDir, "routes.go"), []byte(vendorCode), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &GoEchoDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 0 {
		t.Errorf("got %d endpoints, want 0 (vendor dir should be skipped)", len(endpoints))
	}
}

func TestGoEchoDetectorMultipleMethods(t *testing.T) {
	dir := t.TempDir()

	code := `package main

import "github.com/labstack/echo/v4"

func setup(e *echo.Echo) {
	e.GET("/items", ListItems)
	e.POST("/items", CreateItem)
	e.PUT("/items/:id", UpdateItem)
	e.DELETE("/items/:id", DeleteItem)
	e.HEAD("/items", HeadItems)
	e.OPTIONS("/items", OptionsItems)
}

func ListItems(c echo.Context) error { return nil }
func CreateItem(c echo.Context) error { return nil }
func UpdateItem(c echo.Context) error { return nil }
func DeleteItem(c echo.Context) error { return nil }
func HeadItems(c echo.Context) error { return nil }
func OptionsItems(c echo.Context) error { return nil }
`
	if err := os.WriteFile(filepath.Join(dir, "routes.go"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &GoEchoDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 6 {
		t.Fatalf("got %d endpoints, want 6", len(endpoints))
	}

	// Verify all methods are present.
	methods := make(map[string]bool)
	for _, ep := range endpoints {
		methods[ep.Method] = true
	}
	for _, m := range []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"} {
		if !methods[m] {
			t.Errorf("missing method: %s", m)
		}
	}
}

func TestGoEchoDetectorEmptyDir(t *testing.T) {
	dir := t.TempDir()

	d := &GoEchoDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 0 {
		t.Errorf("got %d endpoints, want 0 for empty dir", len(endpoints))
	}
}

func TestGoEchoDetectorHandlerLocation(t *testing.T) {
	dir := t.TempDir()

	code := `package main

import "github.com/labstack/echo/v4"

func setup(e *echo.Echo) {
	e.GET("/users", ListUsers)
}

func ListUsers(c echo.Context) error {
	return nil
}
`
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &GoEchoDetector{}
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
	if ep.Handler.EndLine == 0 {
		t.Error("handler end line should not be 0")
	}
	if ep.Handler.EndLine < ep.Handler.StartLine {
		t.Errorf("handler end line (%d) < start line (%d)", ep.Handler.EndLine, ep.Handler.StartLine)
	}
}

func TestGoEchoDetectorAddMethod(t *testing.T) {
	dir := t.TempDir()

	code := `package main

import "github.com/labstack/echo/v4"

func setup(e *echo.Echo) {
	e.Add("PATCH", "/users/:id", PatchUser)
	e.Add("PUT", "/users/:id", PutUser)
}

func PatchUser(c echo.Context) error { return nil }
func PutUser(c echo.Context) error { return nil }
`
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &GoEchoDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 2 {
		t.Fatalf("got %d endpoints, want 2", len(endpoints))
	}

	sort.Slice(endpoints, func(i, j int) bool {
		return endpoints[i].Method < endpoints[j].Method
	})

	if endpoints[0].Method != "PATCH" {
		t.Errorf("first endpoint method = %q, want PATCH", endpoints[0].Method)
	}
	if endpoints[0].Name != "PatchUser" {
		t.Errorf("first endpoint name = %q, want PatchUser", endpoints[0].Name)
	}
	if endpoints[1].Method != "PUT" {
		t.Errorf("second endpoint method = %q, want PUT", endpoints[1].Method)
	}
}
