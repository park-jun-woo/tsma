package endpoint

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestTSNextjsDetectorAppRouter(t *testing.T) {
	dir := t.TempDir()

	// Create App Router structure: app/api/users/route.ts
	routeDir := filepath.Join(dir, "app", "api", "users")
	if err := os.MkdirAll(routeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	code := `import { NextRequest, NextResponse } from 'next/server';

export function GET(req: NextRequest) {
  return NextResponse.json([]);
}

export function POST(req: NextRequest) {
  return NextResponse.json({});
}
`
	if err := os.WriteFile(filepath.Join(routeDir, "route.ts"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &TSNextjsDetector{}
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

	// GET endpoint
	if endpoints[0].Method != "GET" {
		t.Errorf("first endpoint method = %q, want GET", endpoints[0].Method)
	}
	if endpoints[0].Path != "/api/users" {
		t.Errorf("first endpoint path = %q, want /api/users", endpoints[0].Path)
	}

	// POST endpoint
	if endpoints[1].Method != "POST" {
		t.Errorf("second endpoint method = %q, want POST", endpoints[1].Method)
	}
	if endpoints[1].Path != "/api/users" {
		t.Errorf("second endpoint path = %q, want /api/users", endpoints[1].Path)
	}

	// Name should be derived from file path: ApiUsers + method
	for _, ep := range endpoints {
		if ep.Status != "todo" {
			t.Errorf("endpoint %s: status = %q, want todo", ep.Name, ep.Status)
		}
	}
}

func TestTSNextjsDetectorPagesRouter(t *testing.T) {
	dir := t.TempDir()

	// Create Pages Router structure: pages/api/users.ts
	pagesDir := filepath.Join(dir, "pages", "api")
	if err := os.MkdirAll(pagesDir, 0o755); err != nil {
		t.Fatal(err)
	}

	code := `import type { NextApiRequest, NextApiResponse } from 'next';

export default function handler(req: NextApiRequest, res: NextApiResponse) {
  res.status(200).json({ name: 'John' });
}
`
	if err := os.WriteFile(filepath.Join(pagesDir, "users.ts"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &TSNextjsDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 1 {
		t.Fatalf("got %d endpoints, want 1", len(endpoints))
	}

	ep := endpoints[0]
	if ep.Method != "ANY" {
		t.Errorf("method = %q, want ANY (Pages Router default export)", ep.Method)
	}
	if ep.Path != "/api/users" {
		t.Errorf("path = %q, want /api/users", ep.Path)
	}
	if ep.Status != "todo" {
		t.Errorf("status = %q, want todo", ep.Status)
	}
}

func TestTSNextjsDetectorBothRouters(t *testing.T) {
	dir := t.TempDir()

	// App Router route
	appDir := filepath.Join(dir, "app", "api", "items")
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		t.Fatal(err)
	}
	appCode := `export async function GET(req) {
  return Response.json([]);
}
`
	if err := os.WriteFile(filepath.Join(appDir, "route.ts"), []byte(appCode), 0o644); err != nil {
		t.Fatal(err)
	}

	// Pages Router route
	pagesDir := filepath.Join(dir, "pages", "api")
	if err := os.MkdirAll(pagesDir, 0o755); err != nil {
		t.Fatal(err)
	}
	pagesCode := `export default function handler(req, res) {
  res.json({});
}
`
	if err := os.WriteFile(filepath.Join(pagesDir, "orders.ts"), []byte(pagesCode), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &TSNextjsDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 2 {
		t.Fatalf("got %d endpoints, want 2 (1 app router + 1 pages router)", len(endpoints))
	}

	// Verify we got both types.
	var appCount, pagesCount int
	for _, ep := range endpoints {
		if ep.Method == "ANY" {
			pagesCount++
		} else {
			appCount++
		}
	}
	if appCount != 1 {
		t.Errorf("app router endpoints = %d, want 1", appCount)
	}
	if pagesCount != 1 {
		t.Errorf("pages router endpoints = %d, want 1", pagesCount)
	}
}

func TestTSNextjsDetectorAppRouterMultipleMethods(t *testing.T) {
	dir := t.TempDir()

	routeDir := filepath.Join(dir, "app", "api", "posts")
	if err := os.MkdirAll(routeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	code := `export function GET(req) {
  return Response.json([]);
}

export async function POST(req) {
  return Response.json({});
}

export function PUT(req) {
  return Response.json({});
}

export function DELETE(req) {
  return Response.json({});
}
`
	if err := os.WriteFile(filepath.Join(routeDir, "route.ts"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &TSNextjsDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 4 {
		t.Fatalf("got %d endpoints, want 4", len(endpoints))
	}

	methods := make(map[string]bool)
	for _, ep := range endpoints {
		methods[ep.Method] = true
		if ep.Path != "/api/posts" {
			t.Errorf("endpoint %s: path = %q, want /api/posts", ep.Name, ep.Path)
		}
	}

	for _, m := range []string{"GET", "POST", "PUT", "DELETE"} {
		if !methods[m] {
			t.Errorf("missing method: %s", m)
		}
	}
}

func TestTSNextjsDetectorSkipsNodeModules(t *testing.T) {
	dir := t.TempDir()

	// Create a valid app route.
	routeDir := filepath.Join(dir, "app", "api", "test")
	if err := os.MkdirAll(routeDir, 0o755); err != nil {
		t.Fatal(err)
	}
	code := `export function GET(req) { return Response.json([]); }
`
	if err := os.WriteFile(filepath.Join(routeDir, "route.ts"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create route inside node_modules (should be skipped).
	nmDir := filepath.Join(dir, "node_modules", "pkg", "app", "api", "hidden")
	if err := os.MkdirAll(nmDir, 0o755); err != nil {
		t.Fatal(err)
	}
	nmCode := `export function GET(req) { return Response.json([]); }
`
	if err := os.WriteFile(filepath.Join(nmDir, "route.ts"), []byte(nmCode), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &TSNextjsDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 1 {
		t.Errorf("got %d endpoints, want 1 (node_modules should be skipped)", len(endpoints))
	}
}

func TestTSNextjsDetectorEmptyDir(t *testing.T) {
	dir := t.TempDir()

	d := &TSNextjsDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 0 {
		t.Errorf("got %d endpoints, want 0", len(endpoints))
	}
}

func TestTSNextjsDetectorEndpointName(t *testing.T) {
	dir := t.TempDir()

	// Test that endpoint names are correctly derived from file paths.
	routeDir := filepath.Join(dir, "app", "api", "users")
	if err := os.MkdirAll(routeDir, 0o755); err != nil {
		t.Fatal(err)
	}
	code := `export function GET(req) { return Response.json([]); }
`
	if err := os.WriteFile(filepath.Join(routeDir, "route.ts"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &TSNextjsDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 1 {
		t.Fatalf("got %d endpoints, want 1", len(endpoints))
	}

	// Name should be derived from path: "ApiUsers" + method "GET" => "ApiUsersGET"
	ep := endpoints[0]
	if ep.Name != "ApiUsersGET" {
		t.Errorf("name = %q, want ApiUsersGET", ep.Name)
	}
}

func TestTSNextjsDetectorPagesRouterIndex(t *testing.T) {
	dir := t.TempDir()

	// Test pages/api/users/index.ts => path should be /api/users
	pagesDir := filepath.Join(dir, "pages", "api", "users")
	if err := os.MkdirAll(pagesDir, 0o755); err != nil {
		t.Fatal(err)
	}
	code := `export default function handler(req, res) {
  res.json([]);
}
`
	if err := os.WriteFile(filepath.Join(pagesDir, "index.ts"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &TSNextjsDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 1 {
		t.Fatalf("got %d endpoints, want 1", len(endpoints))
	}

	ep := endpoints[0]
	if ep.Path != "/api/users" {
		t.Errorf("path = %q, want /api/users", ep.Path)
	}
}

func TestTSNextjsDetectorDynamicSegment(t *testing.T) {
	dir := t.TempDir()

	// Test pages/api/users/[id].ts => path should be /api/users/[id]
	pagesDir := filepath.Join(dir, "pages", "api", "users")
	if err := os.MkdirAll(pagesDir, 0o755); err != nil {
		t.Fatal(err)
	}
	code := `export default function handler(req, res) {
  res.json({});
}
`
	if err := os.WriteFile(filepath.Join(pagesDir, "[id].ts"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &TSNextjsDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 1 {
		t.Fatalf("got %d endpoints, want 1", len(endpoints))
	}

	ep := endpoints[0]
	if ep.Path != "/api/users/[id]" {
		t.Errorf("path = %q, want /api/users/[id]", ep.Path)
	}
}

func TestTSNextjsDetectorSrcPrefix(t *testing.T) {
	dir := t.TempDir()

	// Test src/app/api/users/route.ts => path should still be /api/users
	routeDir := filepath.Join(dir, "src", "app", "api", "users")
	if err := os.MkdirAll(routeDir, 0o755); err != nil {
		t.Fatal(err)
	}
	code := `export function GET(req) { return Response.json([]); }
`
	if err := os.WriteFile(filepath.Join(routeDir, "route.ts"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &TSNextjsDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 1 {
		t.Fatalf("got %d endpoints, want 1", len(endpoints))
	}

	ep := endpoints[0]
	if ep.Path != "/api/users" {
		t.Errorf("path = %q, want /api/users", ep.Path)
	}
}

func TestTSNextjsDetectorHandlerLocation(t *testing.T) {
	dir := t.TempDir()

	routeDir := filepath.Join(dir, "app", "api", "items")
	if err := os.MkdirAll(routeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	code := `export function GET(req) {
  return Response.json([]);
}

export function POST(req) {
  return Response.json({});
}
`
	if err := os.WriteFile(filepath.Join(routeDir, "route.ts"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &TSNextjsDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 2 {
		t.Fatalf("got %d endpoints, want 2", len(endpoints))
	}

	for _, ep := range endpoints {
		if ep.Handler.File == "" {
			t.Errorf("endpoint %s: handler file is empty", ep.Name)
		}
		if ep.Handler.StartLine == 0 {
			t.Errorf("endpoint %s: handler start line should not be 0", ep.Name)
		}
		if ep.Handler.EndLine < ep.Handler.StartLine {
			t.Errorf("endpoint %s: end line (%d) < start line (%d)", ep.Name, ep.Handler.EndLine, ep.Handler.StartLine)
		}
	}
}

func TestTSNextjsDetectorJSRouteFile(t *testing.T) {
	dir := t.TempDir()

	// route.js should also be detected.
	routeDir := filepath.Join(dir, "app", "api", "health")
	if err := os.MkdirAll(routeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	code := `export function GET(req) {
  return Response.json({ status: 'ok' });
}
`
	if err := os.WriteFile(filepath.Join(routeDir, "route.js"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &TSNextjsDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 1 {
		t.Fatalf("got %d endpoints, want 1", len(endpoints))
	}

	if endpoints[0].Method != "GET" {
		t.Errorf("method = %q, want GET", endpoints[0].Method)
	}
}

func TestTSNextjsDetectorPagesRouterName(t *testing.T) {
	dir := t.TempDir()

	pagesDir := filepath.Join(dir, "pages", "api")
	if err := os.MkdirAll(pagesDir, 0o755); err != nil {
		t.Fatal(err)
	}

	code := `export default function handler(req, res) {
  res.json({});
}
`
	if err := os.WriteFile(filepath.Join(pagesDir, "users.ts"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &TSNextjsDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 1 {
		t.Fatalf("got %d endpoints, want 1", len(endpoints))
	}

	// Name should be derived: pages/api/users.ts => "ApiUsers"
	ep := endpoints[0]
	if ep.Name != "ApiUsers" {
		t.Errorf("name = %q, want ApiUsers", ep.Name)
	}
}

func TestDeriveRoutePath(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"app/api/users/route.ts", "/api/users"},
		{"app/api/users/route.js", "/api/users"},
		{"src/app/api/users/route.ts", "/api/users"},
		{"pages/api/users.ts", "/api/users"},
		{"pages/api/users/index.ts", "/api/users"},
		{"pages/api/users/[id].ts", "/api/users/[id]"},
	}

	for _, tt := range tests {
		got := deriveRoutePath(tt.input)
		if got != tt.want {
			t.Errorf("deriveRoutePath(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestDeriveEndpointName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"app/api/users/route.ts", "ApiUsers"},
		{"pages/api/users.ts", "ApiUsers"},
		{"pages/api/users/index.ts", "ApiUsers"},
		{"pages/api/users/[id].ts", "ApiUsersId"},
		{"src/app/api/items/route.ts", "ApiItems"},
	}

	for _, tt := range tests {
		got := deriveEndpointName(tt.input)
		if got != tt.want {
			t.Errorf("deriveEndpointName(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestIsUnderDir(t *testing.T) {
	tests := []struct {
		relPath string
		dir     string
		want    bool
	}{
		{"app/api/users/route.ts", "app", true},
		{"src/app/api/users/route.ts", "app", true},
		{"pages/api/users.ts", "pages/api", true},
		{"src/pages/api/users.ts", "pages/api", true},
		{"other/file.ts", "app", false},
		{"other/file.ts", "pages/api", false},
	}

	for _, tt := range tests {
		got := isUnderDir(tt.relPath, tt.dir)
		if got != tt.want {
			t.Errorf("isUnderDir(%q, %q) = %v, want %v", tt.relPath, tt.dir, got, tt.want)
		}
	}
}

func TestCapitalize(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello", "Hello"},
		{"Hello", "Hello"},
		{"a", "A"},
		{"", ""},
		{"api", "Api"},
	}

	for _, tt := range tests {
		got := capitalize(tt.input)
		if got != tt.want {
			t.Errorf("capitalize(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

// Ensure sort import is used (sort is referenced in TestTSNextjsDetectorAppRouter).
var _ = sort.Slice
