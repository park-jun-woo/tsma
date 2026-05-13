package endpoint

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestPyFastapiDetector(t *testing.T) {
	dir := t.TempDir()

	code := `from fastapi import FastAPI

app = FastAPI()

@app.get("/users")
def list_users():
    pass

@app.post("/users")
def create_user():
    pass
`
	if err := os.WriteFile(filepath.Join(dir, "main.py"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &PyFastapiDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 2 {
		t.Fatalf("got %d endpoints, want 2", len(endpoints))
	}

	expected := map[string]struct {
		method string
		path   string
	}{
		"list_users":  {method: "GET", path: "/users"},
		"create_user": {method: "POST", path: "/users"},
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

func TestPyFastapiDetectorAllMethods(t *testing.T) {
	dir := t.TempDir()

	code := `from fastapi import FastAPI

app = FastAPI()

@app.get("/a")
def handler_get():
    pass

@app.post("/b")
def handler_post():
    pass

@app.put("/c")
def handler_put():
    pass

@app.patch("/d")
def handler_patch():
    pass

@app.delete("/e")
def handler_delete():
    pass

@app.head("/f")
def handler_head():
    pass

@app.options("/g")
def handler_options():
    pass
`
	if err := os.WriteFile(filepath.Join(dir, "routes.py"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &PyFastapiDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 7 {
		t.Fatalf("got %d endpoints, want 7", len(endpoints))
	}

	wantMethods := map[string]string{
		"handler_get":     "GET",
		"handler_post":    "POST",
		"handler_put":     "PUT",
		"handler_patch":   "PATCH",
		"handler_delete":  "DELETE",
		"handler_head":    "HEAD",
		"handler_options": "OPTIONS",
	}

	for _, ep := range endpoints {
		wantMethod, ok := wantMethods[ep.Name]
		if !ok {
			t.Errorf("unexpected endpoint: %s", ep.Name)
			continue
		}
		if ep.Method != wantMethod {
			t.Errorf("endpoint %s: method = %q, want %q", ep.Name, ep.Method, wantMethod)
		}
	}
}

func TestPyFastapiDetectorAsyncHandlers(t *testing.T) {
	dir := t.TempDir()

	code := `from fastapi import FastAPI

app = FastAPI()

@app.get("/items")
async def list_items():
    return []

@app.post("/items")
async def create_item():
    return {}
`
	if err := os.WriteFile(filepath.Join(dir, "main.py"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &PyFastapiDetector{}
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

	if endpoints[0].Name != "create_item" {
		t.Errorf("first endpoint name = %q, want create_item", endpoints[0].Name)
	}
	if endpoints[1].Name != "list_items" {
		t.Errorf("second endpoint name = %q, want list_items", endpoints[1].Name)
	}
}

func TestPyFastapiDetectorRouterInstance(t *testing.T) {
	dir := t.TempDir()

	code := `from fastapi import APIRouter

router = APIRouter()

@router.get("/items")
def get_items():
    return []

@router.post("/items")
def add_item():
    return {}
`
	if err := os.WriteFile(filepath.Join(dir, "routes.py"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &PyFastapiDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 2 {
		t.Fatalf("got %d endpoints, want 2", len(endpoints))
	}
}

func TestPyFastapiDetectorSkipsPycache(t *testing.T) {
	dir := t.TempDir()

	// Create a valid source file.
	code := `from fastapi import FastAPI
app = FastAPI()

@app.get("/hello")
def hello():
    return "hi"
`
	if err := os.WriteFile(filepath.Join(dir, "main.py"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create file inside __pycache__ (should be skipped).
	cacheDir := filepath.Join(dir, "__pycache__")
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		t.Fatal(err)
	}
	cacheCode := `from fastapi import FastAPI
app = FastAPI()

@app.get("/cached")
def cached():
    return "cached"
`
	if err := os.WriteFile(filepath.Join(cacheDir, "main.cpython-310.py"), []byte(cacheCode), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &PyFastapiDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 1 {
		t.Errorf("got %d endpoints, want 1 (__pycache__ should be skipped)", len(endpoints))
	}
}

func TestPyFastapiDetectorSkipsVenv(t *testing.T) {
	dir := t.TempDir()

	venvDir := filepath.Join(dir, "venv", "lib")
	if err := os.MkdirAll(venvDir, 0o755); err != nil {
		t.Fatal(err)
	}
	venvCode := `from fastapi import FastAPI
app = FastAPI()
@app.get("/venv-route")
def venv_handler():
    pass
`
	if err := os.WriteFile(filepath.Join(venvDir, "main.py"), []byte(venvCode), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &PyFastapiDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 0 {
		t.Errorf("got %d endpoints, want 0 (venv should be skipped)", len(endpoints))
	}
}

func TestPyFastapiDetectorEmptyDir(t *testing.T) {
	dir := t.TempDir()

	d := &PyFastapiDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 0 {
		t.Errorf("got %d endpoints, want 0", len(endpoints))
	}
}

func TestPyFastapiDetectorHandlerLocation(t *testing.T) {
	dir := t.TempDir()

	code := `from fastapi import FastAPI

app = FastAPI()

@app.get("/users")
def list_users():
    return []
`
	if err := os.WriteFile(filepath.Join(dir, "main.py"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &PyFastapiDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 1 {
		t.Fatalf("got %d endpoints, want 1", len(endpoints))
	}

	ep := endpoints[0]
	if ep.Handler.File != "main.py" {
		t.Errorf("handler file = %q, want main.py", ep.Handler.File)
	}
	if ep.Handler.StartLine == 0 {
		t.Error("handler start line should not be 0")
	}
	if ep.Handler.EndLine < ep.Handler.StartLine {
		t.Errorf("handler end line (%d) < start line (%d)", ep.Handler.EndLine, ep.Handler.StartLine)
	}
}

func TestPyFastapiDetectorAPIRoute(t *testing.T) {
	dir := t.TempDir()

	code := `from fastapi import FastAPI

app = FastAPI()

@app.api_route("/multi", methods=["GET", "POST"])
def multi_handler():
    pass
`
	if err := os.WriteFile(filepath.Join(dir, "main.py"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &PyFastapiDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 2 {
		t.Fatalf("got %d endpoints, want 2 (one for GET, one for POST)", len(endpoints))
	}

	methods := make(map[string]bool)
	for _, ep := range endpoints {
		methods[ep.Method] = true
		if ep.Name != "multi_handler" {
			t.Errorf("name = %q, want multi_handler", ep.Name)
		}
		if ep.Path != "/multi" {
			t.Errorf("path = %q, want /multi", ep.Path)
		}
	}

	if !methods["GET"] {
		t.Error("missing GET method")
	}
	if !methods["POST"] {
		t.Error("missing POST method")
	}
}

func TestPyFastapiDetectorStackedDecorators(t *testing.T) {
	dir := t.TempDir()

	code := `from fastapi import FastAPI

app = FastAPI()

@app.get("/items")
@app.post("/items")
def item_handler():
    pass
`
	if err := os.WriteFile(filepath.Join(dir, "main.py"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &PyFastapiDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	// Both stacked decorators should produce endpoints.
	if len(endpoints) != 2 {
		t.Fatalf("got %d endpoints, want 2 (stacked decorators)", len(endpoints))
	}

	methods := make(map[string]bool)
	for _, ep := range endpoints {
		methods[ep.Method] = true
	}
	if !methods["GET"] {
		t.Error("missing GET method from stacked decorator")
	}
	if !methods["POST"] {
		t.Error("missing POST method from stacked decorator")
	}
}

func TestPyFastapiDetectorSubdirectory(t *testing.T) {
	dir := t.TempDir()

	subDir := filepath.Join(dir, "routers")
	if err := os.MkdirAll(subDir, 0o755); err != nil {
		t.Fatal(err)
	}

	code := `from fastapi import APIRouter

router = APIRouter()

@router.get("/orders")
def list_orders():
    return []
`
	if err := os.WriteFile(filepath.Join(subDir, "orders.py"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &PyFastapiDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 1 {
		t.Fatalf("got %d endpoints, want 1", len(endpoints))
	}

	if endpoints[0].Name != "list_orders" {
		t.Errorf("name = %q, want list_orders", endpoints[0].Name)
	}
}
