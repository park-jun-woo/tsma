package endpoint

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestTSExpressDetector(t *testing.T) {
	dir := t.TempDir()

	code := `import express from 'express';

const app = express();

function listUsers(req, res) {
  res.json([]);
}

function createUser(req, res) {
  res.json({});
}

function deleteUser(req, res) {
  res.json({});
}

app.get('/users', listUsers);
router.post('/users', createUser);
app.delete('/users/:id', deleteUser);
`
	if err := os.WriteFile(filepath.Join(dir, "app.ts"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &TSExpressDetector{}
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
		"listUsers":  {method: "GET", path: "/users"},
		"createUser": {method: "POST", path: "/users"},
		"deleteUser": {method: "DELETE", path: "/users/:id"},
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

func TestTSExpressDetectorSkipsNodeModules(t *testing.T) {
	dir := t.TempDir()

	// Create a source file in the root.
	code := `import express from 'express';
const app = express();
function hello(req, res) { res.send('hi'); }
app.get('/hello', hello);
`
	if err := os.WriteFile(filepath.Join(dir, "index.ts"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create a file inside node_modules that should be skipped.
	nmDir := filepath.Join(dir, "node_modules", "some-pkg")
	if err := os.MkdirAll(nmDir, 0o755); err != nil {
		t.Fatal(err)
	}
	nmCode := `const router = require('express').Router();
function internal(req, res) { res.send('internal'); }
router.get('/internal', internal);
`
	if err := os.WriteFile(filepath.Join(nmDir, "index.js"), []byte(nmCode), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &TSExpressDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 1 {
		t.Errorf("got %d endpoints, want 1 (node_modules should be skipped)", len(endpoints))
	}
	if len(endpoints) > 0 && endpoints[0].Name != "hello" {
		t.Errorf("endpoint name = %q, want hello", endpoints[0].Name)
	}
}

func TestTSExpressDetectorMethodNormalization(t *testing.T) {
	dir := t.TempDir()

	// Express uses lowercase method names (get, post, etc.);
	// the detector should normalize them to uppercase.
	code := `import express from 'express';
const app = express();

function getHandler(req, res) {}
function postHandler(req, res) {}
function putHandler(req, res) {}
function patchHandler(req, res) {}
function deleteHandler(req, res) {}

app.get('/a', getHandler);
app.post('/b', postHandler);
app.put('/c', putHandler);
app.patch('/d', patchHandler);
app.delete('/e', deleteHandler);
`
	if err := os.WriteFile(filepath.Join(dir, "routes.ts"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &TSExpressDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 5 {
		t.Fatalf("got %d endpoints, want 5", len(endpoints))
	}

	wantMethods := map[string]string{
		"getHandler":    "GET",
		"postHandler":   "POST",
		"putHandler":    "PUT",
		"patchHandler":  "PATCH",
		"deleteHandler": "DELETE",
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

func TestTSExpressDetectorJSFiles(t *testing.T) {
	dir := t.TempDir()

	code := `const express = require('express');
const app = express();
function healthCheck(req, res) { res.send('ok'); }
app.get('/health', healthCheck);
`
	if err := os.WriteFile(filepath.Join(dir, "app.js"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &TSExpressDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 1 {
		t.Fatalf("got %d endpoints, want 1", len(endpoints))
	}

	if endpoints[0].Name != "healthCheck" {
		t.Errorf("name = %q, want healthCheck", endpoints[0].Name)
	}
}

func TestTSExpressDetectorSkipsDTS(t *testing.T) {
	dir := t.TempDir()

	// .d.ts files should be skipped.
	dtsCode := `declare function listUsers(req: any, res: any): void;
app.get('/users', listUsers);
`
	if err := os.WriteFile(filepath.Join(dir, "types.d.ts"), []byte(dtsCode), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &TSExpressDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 0 {
		t.Errorf("got %d endpoints, want 0 (.d.ts files should be skipped)", len(endpoints))
	}
}

func TestTSExpressDetectorEmptyDir(t *testing.T) {
	dir := t.TempDir()

	d := &TSExpressDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 0 {
		t.Errorf("got %d endpoints, want 0", len(endpoints))
	}
}

func TestTSExpressDetectorRouterPattern(t *testing.T) {
	dir := t.TempDir()

	code := `import { Router } from 'express';
const router = Router();

function getItems(req, res) { res.json([]); }
function addItem(req, res) { res.json({}); }

router.get('/items', getItems);
router.post('/items', addItem);
`
	if err := os.WriteFile(filepath.Join(dir, "items.ts"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &TSExpressDetector{}
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

	if endpoints[0].Name != "addItem" {
		t.Errorf("first endpoint = %q, want addItem", endpoints[0].Name)
	}
	if endpoints[1].Name != "getItems" {
		t.Errorf("second endpoint = %q, want getItems", endpoints[1].Name)
	}
}

func TestTSExpressDetectorHandlerLocation(t *testing.T) {
	dir := t.TempDir()

	code := `import express from 'express';
const app = express();

function listUsers(req, res) {
  res.json([]);
}

app.get('/users', listUsers);
`
	if err := os.WriteFile(filepath.Join(dir, "app.ts"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}

	d := &TSExpressDetector{}
	endpoints, err := d.Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	if len(endpoints) != 1 {
		t.Fatalf("got %d endpoints, want 1", len(endpoints))
	}

	ep := endpoints[0]
	if ep.Handler.File != "app.ts" {
		t.Errorf("handler file = %q, want app.ts", ep.Handler.File)
	}
	if ep.Handler.StartLine == 0 {
		t.Error("handler start line should not be 0")
	}
}
