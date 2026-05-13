package graph

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// ---------------------------------------------------------------------------
// Helper: write file into dir, creating intermediate directories.
// ---------------------------------------------------------------------------

func writeFile(t *testing.T, dir, rel, content string) {
	t.Helper()
	abs := filepath.Join(dir, rel)
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(abs, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// ---------------------------------------------------------------------------
// Factory tests
// ---------------------------------------------------------------------------

func TestNewBuilderGo(t *testing.T) {
	b := NewBuilder("go")
	if _, ok := b.(*GoBuilder); !ok {
		t.Errorf("NewBuilder(\"go\") returned %T, want *GoBuilder", b)
	}
}

func TestNewBuilderTypescript(t *testing.T) {
	b := NewBuilder("typescript")
	if _, ok := b.(*TSBuilder); !ok {
		t.Errorf("NewBuilder(\"typescript\") returned %T, want *TSBuilder", b)
	}
}

func TestNewBuilderPython(t *testing.T) {
	b := NewBuilder("python")
	if _, ok := b.(*PyBuilder); !ok {
		t.Errorf("NewBuilder(\"python\") returned %T, want *PyBuilder", b)
	}
}

func TestNewBuilderUnsupported(t *testing.T) {
	b := NewBuilder("ruby")
	u, ok := b.(*UnsupportedBuilder)
	if !ok {
		t.Fatalf("NewBuilder(\"ruby\") returned %T, want *UnsupportedBuilder", b)
	}
	if u.Lang != "ruby" {
		t.Errorf("UnsupportedBuilder.Lang = %q, want \"ruby\"", u.Lang)
	}
	_, _, err := b.Build(".", nil)
	if err == nil {
		t.Fatal("expected error from UnsupportedBuilder.Build")
	}
}

// ---------------------------------------------------------------------------
// GoBuilder tests
// ---------------------------------------------------------------------------

func TestGoBuilderSamePackageCall(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "main.go", `package main

func main() {
	handleLogin()
}

func handleLogin() {
	generateToken()
}

func generateToken() string {
	return "token"
}
`)

	functions := []model.Function{
		{QualifiedName: "main", Name: "main", File: "main.go", StartLine: 3, EndLine: 5, Exported: false, Status: model.StatusTodo},
		{QualifiedName: "handleLogin", Name: "handleLogin", File: "main.go", StartLine: 7, EndLine: 9, Exported: false, Status: model.StatusTodo},
		{QualifiedName: "generateToken", Name: "generateToken", File: "main.go", StartLine: 11, EndLine: 13, Exported: false, Status: model.StatusTodo},
	}

	b := &GoBuilder{}
	result, summary, err := b.Build(dir, functions)
	if err != nil {
		t.Fatalf("GoBuilder.Build: %v", err)
	}

	if summary.Nodes != 3 {
		t.Errorf("summary.Nodes = %d, want 3", summary.Nodes)
	}

	// main calls handleLogin.
	mainFn := findFunc(result, "main")
	if mainFn == nil {
		t.Fatal("expected to find 'main'")
	}
	if !mainFn.EntryPoint {
		t.Error("main should be an entry point")
	}
	if len(mainFn.Callees) == 0 {
		t.Error("main should have callees")
	} else {
		found := false
		for _, e := range mainFn.Callees {
			if e.Target == "handleLogin" {
				found = true
			}
		}
		if !found {
			t.Error("main should call handleLogin")
		}
	}

	// handleLogin calls generateToken.
	hlFn := findFunc(result, "handleLogin")
	if hlFn == nil {
		t.Fatal("expected to find 'handleLogin'")
	}
	if len(hlFn.Callees) == 0 {
		t.Error("handleLogin should have callees")
	}

	// handleLogin should have main as a caller.
	if len(hlFn.Callers) == 0 {
		t.Error("handleLogin should have callers")
	}

	// generateToken has no callees and is not dead (has callers).
	gtFn := findFunc(result, "generateToken")
	if gtFn == nil {
		t.Fatal("expected to find 'generateToken'")
	}
	if len(gtFn.Callers) == 0 {
		t.Error("generateToken should have callers")
	}
	if gtFn.Dead {
		t.Error("generateToken should not be dead (has callers)")
	}
}

func TestGoBuilderImportCall(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "main.go", `package main

import "example/internal/auth"

func main() {
	auth.Login()
}
`)

	writeFile(t, dir, "internal/auth/handler.go", `package auth

func Login() error {
	return nil
}
`)

	functions := []model.Function{
		{QualifiedName: "main", Name: "main", File: "main.go", StartLine: 5, EndLine: 7, Exported: false, Status: model.StatusTodo},
		{QualifiedName: "internal/auth.Login", Name: "Login", File: "internal/auth/handler.go", StartLine: 3, EndLine: 5, Exported: true, Status: model.StatusTodo},
	}

	b := &GoBuilder{}
	result, _, err := b.Build(dir, functions)
	if err != nil {
		t.Fatalf("GoBuilder.Build: %v", err)
	}

	mainFn := findFunc(result, "main")
	if mainFn == nil {
		t.Fatal("expected to find 'main'")
	}

	// main should call Login via import resolution.
	found := false
	for _, e := range mainFn.Callees {
		if e.Target == "internal/auth.Login" {
			found = true
			if e.Ambiguous {
				t.Error("import-based call should not be ambiguous")
			}
		}
	}
	if !found {
		t.Error("main should call internal/auth.Login")
	}
}

func TestGoBuilderMethodAmbiguous(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "main.go", `package main

func main() {
	h.Login()
}
`)

	writeFile(t, dir, "auth.go", `package main

type AuthHandler struct{}

func (a *AuthHandler) Login() error {
	return nil
}
`)

	writeFile(t, dir, "user.go", `package main

type UserHandler struct{}

func (u *UserHandler) Login() error {
	return nil
}
`)

	functions := []model.Function{
		{QualifiedName: "main", Name: "main", File: "main.go", StartLine: 3, EndLine: 5, Exported: false, Status: model.StatusTodo},
		{QualifiedName: "AuthHandler.Login", Name: "Login", File: "auth.go", StartLine: 5, EndLine: 7, IsMethod: true, Receiver: "AuthHandler", Exported: true, Status: model.StatusTodo},
		{QualifiedName: "UserHandler.Login", Name: "Login", File: "user.go", StartLine: 5, EndLine: 7, IsMethod: true, Receiver: "UserHandler", Exported: true, Status: model.StatusTodo},
	}

	b := &GoBuilder{}
	result, _, err := b.Build(dir, functions)
	if err != nil {
		t.Fatalf("GoBuilder.Build: %v", err)
	}

	mainFn := findFunc(result, "main")
	if mainFn == nil {
		t.Fatal("expected to find 'main'")
	}

	// main should call both Login methods with ambiguous=true.
	if len(mainFn.Callees) < 2 {
		t.Fatalf("main should have at least 2 callees, got %d", len(mainFn.Callees))
	}

	for _, e := range mainFn.Callees {
		if !e.Ambiguous {
			t.Errorf("edge to %q should be ambiguous", e.Target)
		}
	}
}

func TestGoBuilderDeadFunction(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "main.go", `package main

func main() {}

func unused() {}
`)

	functions := []model.Function{
		{QualifiedName: "main", Name: "main", File: "main.go", StartLine: 3, EndLine: 3, Exported: false, Status: model.StatusTodo},
		{QualifiedName: "unused", Name: "unused", File: "main.go", StartLine: 5, EndLine: 5, Exported: false, Status: model.StatusTodo},
	}

	b := &GoBuilder{}
	result, summary, err := b.Build(dir, functions)
	if err != nil {
		t.Fatalf("GoBuilder.Build: %v", err)
	}

	mainFn := findFunc(result, "main")
	if mainFn == nil || !mainFn.EntryPoint {
		t.Error("main should be an entry point")
	}

	unusedFn := findFunc(result, "unused")
	if unusedFn == nil {
		t.Fatal("expected to find 'unused'")
	}
	if !unusedFn.Dead {
		t.Error("unused should be marked as dead")
	}

	if summary.Dead != 1 {
		t.Errorf("summary.Dead = %d, want 1", summary.Dead)
	}
	if summary.EntryPoints != 1 {
		t.Errorf("summary.EntryPoints = %d, want 1", summary.EntryPoints)
	}
}

func TestGoBuilderExportedEntryPoint(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "handler.go", `package api

func HandleLogin() {}

func helperFunc() {}
`)

	functions := []model.Function{
		{QualifiedName: "HandleLogin", Name: "HandleLogin", File: "handler.go", StartLine: 3, EndLine: 3, Exported: true, Status: model.StatusTodo},
		{QualifiedName: "helperFunc", Name: "helperFunc", File: "handler.go", StartLine: 5, EndLine: 5, Exported: false, Status: model.StatusTodo},
	}

	b := &GoBuilder{}
	result, _, err := b.Build(dir, functions)
	if err != nil {
		t.Fatalf("GoBuilder.Build: %v", err)
	}

	hl := findFunc(result, "HandleLogin")
	if hl == nil || !hl.EntryPoint {
		t.Error("HandleLogin should be an entry point (exported)")
	}

	helper := findFunc(result, "helperFunc")
	if helper == nil || !helper.Dead {
		t.Error("helperFunc should be dead (unexported, no callers)")
	}
}

func TestGoBuilderEmptyFunctions(t *testing.T) {
	dir := t.TempDir()

	b := &GoBuilder{}
	result, summary, err := b.Build(dir, nil)
	if err != nil {
		t.Fatalf("GoBuilder.Build: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 functions, got %d", len(result))
	}
	if summary.Nodes != 0 {
		t.Errorf("summary.Nodes = %d, want 0", summary.Nodes)
	}
}

// ---------------------------------------------------------------------------
// TSBuilder tests
// ---------------------------------------------------------------------------

func TestTSBuilderBasic(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "handler.ts", `import { login } from './auth-service';

export async function handleLogin(req: Request) {
  const result = await login(req.body);
  return result;
}
`)

	writeFile(t, dir, "auth-service.ts", `import { generateToken } from './utils';

export async function login(credentials: any) {
  const token = generateToken(credentials);
  return token;
}
`)

	writeFile(t, dir, "utils.ts", `export function generateToken(user: any): string {
  return 'jwt-token-' + user.id;
}
`)

	functions := []model.Function{
		{QualifiedName: "handleLogin", Name: "handleLogin", File: "handler.ts", StartLine: 3, EndLine: 6, Exported: true, Status: model.StatusTodo},
		{QualifiedName: "login", Name: "login", File: "auth-service.ts", StartLine: 3, EndLine: 6, Exported: true, Status: model.StatusTodo},
		{QualifiedName: "generateToken", Name: "generateToken", File: "utils.ts", StartLine: 1, EndLine: 3, Exported: true, Status: model.StatusTodo},
	}

	b := &TSBuilder{}
	result, summary, err := b.Build(dir, functions)
	if err != nil {
		t.Fatalf("TSBuilder.Build: %v", err)
	}

	if summary.Nodes != 3 {
		t.Errorf("summary.Nodes = %d, want 3", summary.Nodes)
	}

	// handleLogin should call login.
	hlFn := findFunc(result, "handleLogin")
	if hlFn == nil {
		t.Fatal("expected to find 'handleLogin'")
	}
	foundLogin := false
	for _, e := range hlFn.Callees {
		if e.Target == "login" {
			foundLogin = true
		}
	}
	if !foundLogin {
		t.Error("handleLogin should call login")
	}

	// login should call generateToken.
	loginFn := findFunc(result, "login")
	if loginFn == nil {
		t.Fatal("expected to find 'login'")
	}
	foundToken := false
	for _, e := range loginFn.Callees {
		if e.Target == "generateToken" {
			foundToken = true
		}
	}
	if !foundToken {
		t.Error("login should call generateToken")
	}
}

func TestTSBuilderSkipsBuiltins(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "handler.ts", `export function handler() {
  console.log("hello");
  JSON.parse("{}");
  Math.random();
}
`)

	functions := []model.Function{
		{QualifiedName: "handler", Name: "handler", File: "handler.ts", StartLine: 1, EndLine: 5, Exported: true, Status: model.StatusTodo},
	}

	b := &TSBuilder{}
	result, _, err := b.Build(dir, functions)
	if err != nil {
		t.Fatalf("TSBuilder.Build: %v", err)
	}

	hlFn := findFunc(result, "handler")
	if hlFn == nil {
		t.Fatal("expected to find 'handler'")
	}

	// Should have no callees (all built-in).
	for _, e := range hlFn.Callees {
		t.Errorf("unexpected callee: %q (built-in calls should be skipped)", e.Target)
	}
}

func TestTSBuilderDeadFunction(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "main.ts", `export function main() {}

function unused() {}
`)

	functions := []model.Function{
		{QualifiedName: "main", Name: "main", File: "main.ts", StartLine: 1, EndLine: 1, Exported: true, Status: model.StatusTodo},
		{QualifiedName: "unused", Name: "unused", File: "main.ts", StartLine: 3, EndLine: 3, Exported: false, Status: model.StatusTodo},
	}

	b := &TSBuilder{}
	result, summary, err := b.Build(dir, functions)
	if err != nil {
		t.Fatalf("TSBuilder.Build: %v", err)
	}

	mainFn := findFunc(result, "main")
	if mainFn == nil || !mainFn.EntryPoint {
		t.Error("main (exported) should be an entry point")
	}

	unusedFn := findFunc(result, "unused")
	if unusedFn == nil || !unusedFn.Dead {
		t.Error("unused should be dead")
	}

	if summary.Dead != 1 {
		t.Errorf("summary.Dead = %d, want 1", summary.Dead)
	}
}

// ---------------------------------------------------------------------------
// PyBuilder tests
// ---------------------------------------------------------------------------

func TestPyBuilderBasic(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "handler.py", `from auth_service import login

def handle_login(request):
    result = login(request.body)
    return result
`)

	writeFile(t, dir, "auth_service.py", `from utils import generate_token

def login(credentials):
    token = generate_token(credentials)
    return token
`)

	writeFile(t, dir, "utils.py", `def generate_token(user):
    return "jwt-token-" + str(user.id)
`)

	functions := []model.Function{
		{QualifiedName: "handle_login", Name: "handle_login", File: "handler.py", StartLine: 3, EndLine: 5, Exported: false, Status: model.StatusTodo},
		{QualifiedName: "login", Name: "login", File: "auth_service.py", StartLine: 3, EndLine: 5, Exported: false, Status: model.StatusTodo},
		{QualifiedName: "generate_token", Name: "generate_token", File: "utils.py", StartLine: 1, EndLine: 2, Exported: false, Status: model.StatusTodo},
	}

	b := &PyBuilder{}
	result, summary, err := b.Build(dir, functions)
	if err != nil {
		t.Fatalf("PyBuilder.Build: %v", err)
	}

	if summary.Nodes != 3 {
		t.Errorf("summary.Nodes = %d, want 3", summary.Nodes)
	}

	// handle_login should call login.
	hlFn := findFunc(result, "handle_login")
	if hlFn == nil {
		t.Fatal("expected to find 'handle_login'")
	}
	foundLogin := false
	for _, e := range hlFn.Callees {
		if e.Target == "login" {
			foundLogin = true
		}
	}
	if !foundLogin {
		t.Error("handle_login should call login")
	}

	// login should call generate_token.
	loginFn := findFunc(result, "login")
	if loginFn == nil {
		t.Fatal("expected to find 'login'")
	}
	foundToken := false
	for _, e := range loginFn.Callees {
		if e.Target == "generate_token" {
			foundToken = true
		}
	}
	if !foundToken {
		t.Error("login should call generate_token")
	}
}

func TestPyBuilderSkipsBuiltins(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "handler.py", `def handler():
    print("hello")
    x = len([1, 2, 3])
    y = str(42)
`)

	functions := []model.Function{
		{QualifiedName: "handler", Name: "handler", File: "handler.py", StartLine: 1, EndLine: 4, Exported: false, Status: model.StatusTodo},
	}

	b := &PyBuilder{}
	result, _, err := b.Build(dir, functions)
	if err != nil {
		t.Fatalf("PyBuilder.Build: %v", err)
	}

	hlFn := findFunc(result, "handler")
	if hlFn == nil {
		t.Fatal("expected to find 'handler'")
	}

	for _, e := range hlFn.Callees {
		t.Errorf("unexpected callee: %q (built-in calls should be skipped)", e.Target)
	}
}

func TestPyBuilderDeadFunction(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "main.py", `def main():
    pass

def unused():
    pass
`)

	functions := []model.Function{
		{QualifiedName: "main", Name: "main", File: "main.py", StartLine: 1, EndLine: 2, Exported: false, Status: model.StatusTodo},
		{QualifiedName: "unused", Name: "unused", File: "main.py", StartLine: 4, EndLine: 5, Exported: false, Status: model.StatusTodo},
	}

	b := &PyBuilder{}
	result, summary, err := b.Build(dir, functions)
	if err != nil {
		t.Fatalf("PyBuilder.Build: %v", err)
	}

	mainFn := findFunc(result, "main")
	if mainFn == nil || !mainFn.EntryPoint {
		t.Error("main should be an entry point")
	}

	unusedFn := findFunc(result, "unused")
	if unusedFn == nil || !unusedFn.Dead {
		t.Error("unused should be dead")
	}

	if summary.Dead != 1 {
		t.Errorf("summary.Dead = %d, want 1", summary.Dead)
	}
}

// ---------------------------------------------------------------------------
// buildSummary tests
// ---------------------------------------------------------------------------

func TestBuildSummary(t *testing.T) {
	functions := []model.Function{
		{QualifiedName: "main", Name: "main", EntryPoint: true, Callees: []model.Edge{{Target: "a"}, {Target: "b"}}},
		{QualifiedName: "a", Name: "a", Callees: []model.Edge{{Target: "c"}}},
		{QualifiedName: "b", Name: "b", Dead: true},
		{QualifiedName: "c", Name: "c"},
	}

	s := buildSummary(functions)
	if s.Nodes != 4 {
		t.Errorf("Nodes = %d, want 4", s.Nodes)
	}
	if s.Edges != 3 {
		t.Errorf("Edges = %d, want 3", s.Edges)
	}
	if s.EntryPoints != 1 {
		t.Errorf("EntryPoints = %d, want 1", s.EntryPoints)
	}
	if s.Dead != 1 {
		t.Errorf("Dead = %d, want 1", s.Dead)
	}
}

// ---------------------------------------------------------------------------
// markEntryAndDead tests
// ---------------------------------------------------------------------------

func TestMarkEntryAndDead(t *testing.T) {
	functions := []model.Function{
		{Name: "main", Callers: nil, Exported: false},
		{Name: "init", Callers: nil, Exported: false},
		{Name: "TestFoo", Callers: nil, Exported: true},
		{Name: "HandleLogin", Callers: nil, Exported: true},
		{Name: "helperFunc", Callers: nil, Exported: false},
		{Name: "usedFunc", Callers: []model.Edge{{Target: "main"}}, Exported: false},
	}

	markEntryAndDead(functions, true)

	if !functions[0].EntryPoint {
		t.Error("main should be entry point")
	}
	if !functions[1].EntryPoint {
		t.Error("init should be entry point")
	}
	if !functions[2].EntryPoint {
		t.Error("TestFoo should be entry point")
	}
	if !functions[3].EntryPoint {
		t.Error("HandleLogin (exported) should be entry point")
	}
	if !functions[4].Dead {
		t.Error("helperFunc should be dead")
	}
	if functions[5].Dead || functions[5].EntryPoint {
		t.Error("usedFunc should be neither dead nor entry point")
	}
}

// ---------------------------------------------------------------------------
// Helper
// ---------------------------------------------------------------------------

func findFunc(functions []model.Function, name string) *model.Function {
	for i := range functions {
		if functions[i].Name == name {
			return &functions[i]
		}
	}
	return nil
}
