package golang

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/neostd/go/xscript/internal/testutil"
)

func init() {
	testutil.SetupMise(NAME, NAME)
}

func TestNew(t *testing.T) {
	testutil.RequireExecutable(t, NAME)

	cmd := New("arg1", "arg2")
	expectedArgs := []string{"arg1", "arg2"}
	testutil.AssertExecutablePath(t, cmd.Path, NAME)
	if !reflect.DeepEqual(cmd.Args[1:], expectedArgs) {
		t.Fatalf("expected args %v, got %v", expectedArgs, cmd.Args[1:])
	}
}

func TestFile(t *testing.T) {
	testutil.RequireExecutable(t, NAME)

	cmd := File("test.go", "arg1")
	expectedArgs := []string{"run", "test.go", "arg1"}
	testutil.AssertExecutablePath(t, cmd.Path, NAME)
	if !reflect.DeepEqual(cmd.Args[1:], expectedArgs) {
		t.Fatalf("expected args %v, got %v", expectedArgs, cmd.Args[1:])
	}
}

func TestInline(t *testing.T) {
	testutil.RequireExecutable(t, NAME)

	cmd := Inline("package main\nfunc main() {}", "arg1")
	testutil.AssertExecutablePath(t, cmd.Path, NAME)
	if cmd.TempFile == nil {
		t.Fatal("expected temp file to be set")
	}
	if filepath.Ext(*cmd.TempFile) != ".go" {
		t.Fatalf("expected temp file extension .go, got %q", *cmd.TempFile)
	}
	if !strings.HasPrefix(filepath.Base(*cmd.TempFile), "run-") {
		t.Fatalf("expected temp file name to start with run-, got %q", *cmd.TempFile)
	}
	if !reflect.DeepEqual(cmd.Args[1:3], []string{"run", *cmd.TempFile}) {
		t.Fatalf("expected args prefix %v, got %v", []string{"run", *cmd.TempFile}, cmd.Args[1:3])
	}
	if !reflect.DeepEqual(cmd.Args[3:], []string{"arg1"}) {
		t.Fatalf("expected trailing args [arg1], got %v", cmd.Args[3:])
	}
}

func TestScriptFile(t *testing.T) {
	testutil.RequireExecutable(t, NAME)

	cmd := Script("test.go", "arg1")
	expectedArgs := []string{"run", "test.go", "arg1"}
	testutil.AssertExecutablePath(t, cmd.Path, NAME)
	if !reflect.DeepEqual(cmd.Args[1:], expectedArgs) {
		t.Fatalf("expected args %v, got %v", expectedArgs, cmd.Args[1:])
	}
}

func TestScriptInline(t *testing.T) {
	testutil.RequireExecutable(t, NAME)

	cmd := Script("package main\nfunc main() {}", "arg1")
	testutil.AssertExecutablePath(t, cmd.Path, NAME)
	if cmd.TempFile == nil {
		t.Fatal("expected temp file to be set")
	}
	if filepath.Ext(*cmd.TempFile) != ".go" {
		t.Fatalf("expected temp file extension .go, got %q", *cmd.TempFile)
	}
	if !reflect.DeepEqual(cmd.Args[1:3], []string{"run", *cmd.TempFile}) {
		t.Fatalf("expected args prefix %v, got %v", []string{"run", *cmd.TempFile}, cmd.Args[1:3])
	}
	if !reflect.DeepEqual(cmd.Args[3:], []string{"arg1"}) {
		t.Fatalf("expected trailing args [arg1], got %v", cmd.Args[3:])
	}
}
