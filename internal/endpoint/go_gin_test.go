package endpoint

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGoGinDetector(t *testing.T) {
	dir := t.TempDir()

	// Create a Go source file with Gin routes.
	code := `package main

import "github.com/gin-gonic/gin"

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/users", ListUsers)
	r.POST("/users", CreateUser)
	r.GET("/users/:id", GetUser)
	r.PUT("/users/:id", UpdateUser)
	r.DELETE("/users/:id", DeleteUser)
	return r
}

func ListUsers(c *gin.Context) {}
func CreateUser(c *gin.Context) {}
func GetUser(c *gin.Context) {}
func UpdateUser(c *gin.Context) {}
func DeleteUser(c *gin.Context) {}
`
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &GoGinDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 5 {
		t.Fatalf("got %d endpoints, want 5", len(endpoints))
	}

	expected := map[string]string{
		"ListUsers":  "GET",
		"CreateUser": "POST",
		"GetUser":    "GET",
		"UpdateUser": "PUT",
		"DeleteUser": "DELETE",
	}

	for _, ep := range endpoints {
		wantMethod, ok := expected[ep.Name]
		if !ok {
			t.Errorf("unexpected endpoint: %s", ep.Name)
			continue
		}
		if ep.Method != wantMethod {
			t.Errorf("endpoint %s: method = %q, want %q", ep.Name, ep.Method, wantMethod)
		}
		if ep.Status != "todo" {
			t.Errorf("endpoint %s: status = %q, want todo", ep.Name, ep.Status)
		}
	}
}

func TestGoGinDetectorSkipsTestFiles(t *testing.T) {
	dir := t.TempDir()

	// Regular source file.
	code := `package main
import "github.com/gin-gonic/gin"
func setup(r *gin.Engine) {
	r.GET("/hello", Hello)
}
func Hello(c *gin.Context) {}
`
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	// Test file — should be skipped.
	testCode := `package main
import "github.com/gin-gonic/gin"
func setupTest(r *gin.Engine) {
	r.GET("/test", TestHandler)
}
func TestHandler(c *gin.Context) {}
`
	if err := os.WriteFile(filepath.Join(dir, "main_test.go"), []byte(testCode), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &GoGinDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 1 {
		t.Errorf("got %d endpoints, want 1 (test file should be skipped)", len(endpoints))
	}
}
