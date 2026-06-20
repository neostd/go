package dotnet

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

	cmd := File("test.cs", "arg1")
	expectedArgs := []string{"run", "--file", "test.cs", "--", "arg1"}
	testutil.AssertExecutablePath(t, cmd.Path, NAME)
	if !reflect.DeepEqual(cmd.Args[1:], expectedArgs) {
		t.Fatalf("expected args %v, got %v", expectedArgs, cmd.Args[1:])
	}
}

func TestInline(t *testing.T) {
	testutil.RequireExecutable(t, NAME)

	cmd := Inline("Console.WriteLine(\"hello\");", "arg1")
	testutil.AssertExecutablePath(t, cmd.Path, NAME)
	if cmd.TempFile == nil {
		t.Fatal("expected temp file to be set")
	}
	if filepath.Ext(*cmd.TempFile) != ".cs" {
		t.Fatalf("expected temp file extension .cs, got %q", *cmd.TempFile)
	}
	if !strings.HasPrefix(filepath.Base(*cmd.TempFile), "run-") {
		t.Fatalf("expected temp file name to start with run-, got %q", *cmd.TempFile)
	}
	expectedPrefix := []string{"run", "--file"}
	if !reflect.DeepEqual(cmd.Args[1:3], expectedPrefix) {
		t.Fatalf("expected args prefix %v, got %v", expectedPrefix, cmd.Args[1:3])
	}
	if cmd.Args[3] != *cmd.TempFile {
		t.Fatalf("expected temp file arg %q, got %q", *cmd.TempFile, cmd.Args[3])
	}
	if !reflect.DeepEqual(cmd.Args[4:], []string{"--", "arg1"}) {
		t.Fatalf("expected trailing args %v, got %v", []string{"--", "arg1"}, cmd.Args[4:])
	}
}

func TestScriptFile(t *testing.T) {
	testutil.RequireExecutable(t, NAME)

	cmd := Script("test.cs", "arg1")
	expectedArgs := []string{"run", "--file", "test.cs", "--", "arg1"}
	testutil.AssertExecutablePath(t, cmd.Path, NAME)
	if !reflect.DeepEqual(cmd.Args[1:], expectedArgs) {
		t.Fatalf("expected args %v, got %v", expectedArgs, cmd.Args[1:])
	}
}

func TestScriptInline(t *testing.T) {
	testutil.RequireExecutable(t, NAME)

	cmd := Script("Console.WriteLine(\"hello\");", "arg1")
	testutil.AssertExecutablePath(t, cmd.Path, NAME)
	if cmd.TempFile == nil {
		t.Fatal("expected temp file to be set")
	}
	if filepath.Ext(*cmd.TempFile) != ".cs" {
		t.Fatalf("expected temp file extension .cs, got %q", *cmd.TempFile)
	}
	if !reflect.DeepEqual(cmd.Args[1:3], []string{"run", "--file"}) {
		t.Fatalf("expected args prefix %v, got %v", []string{"run", "--file"}, cmd.Args[1:3])
	}
	if !reflect.DeepEqual(cmd.Args[4:], []string{"--", "arg1"}) {
		t.Fatalf("expected trailing args %v, got %v", []string{"--", "arg1"}, cmd.Args[4:])
	}
}
