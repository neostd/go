package bash

import (
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

	cmd := File("test.sh", "arg1")
	testutil.AssertExecutablePath(t, cmd.Path, NAME)
	if !reflect.DeepEqual(cmd.Args[1:1+len(ScriptArgs)], ScriptArgs) {
		t.Fatalf("expected script args %v, got %v", ScriptArgs, cmd.Args[1:1+len(ScriptArgs)])
	}
	if !strings.HasSuffix(cmd.Args[1+len(ScriptArgs)], "test.sh") {
		t.Fatalf("expected script file ending with test.sh, got %q", cmd.Args[1+len(ScriptArgs)])
	}
	if !reflect.DeepEqual(cmd.Args[2+len(ScriptArgs):], []string{"arg1"}) {
		t.Fatalf("expected trailing args [arg1], got %v", cmd.Args[2+len(ScriptArgs):])
	}
}

func TestInline(t *testing.T) {
	testutil.RequireExecutable(t, NAME)

	cmd := Inline("echo hello", "arg1")
	expectedArgs := append(append([]string{}, ScriptArgs...), "-c", "echo hello", "arg1")
	testutil.AssertExecutablePath(t, cmd.Path, NAME)
	if !reflect.DeepEqual(cmd.Args[1:], expectedArgs) {
		t.Fatalf("expected args %v, got %v", expectedArgs, cmd.Args[1:])
	}
}

func TestScriptFile(t *testing.T) {
	testutil.RequireExecutable(t, NAME)

	cmd := Script("test.sh", "arg1")
	testutil.AssertExecutablePath(t, cmd.Path, NAME)
	if !reflect.DeepEqual(cmd.Args[1:1+len(ScriptArgs)], ScriptArgs) {
		t.Fatalf("expected script args %v, got %v", ScriptArgs, cmd.Args[1:1+len(ScriptArgs)])
	}
	if !strings.HasSuffix(cmd.Args[1+len(ScriptArgs)], "test.sh") {
		t.Fatalf("expected script file ending with test.sh, got %q", cmd.Args[1+len(ScriptArgs)])
	}
	if !reflect.DeepEqual(cmd.Args[2+len(ScriptArgs):], []string{"arg1"}) {
		t.Fatalf("expected trailing args [arg1], got %v", cmd.Args[2+len(ScriptArgs):])
	}
}

func TestScriptInline(t *testing.T) {
	testutil.RequireExecutable(t, NAME)

	cmd := Script("echo hello", "arg1")
	expectedArgs := append(append([]string{}, ScriptArgs...), "-c", "echo hello", "arg1")
	testutil.AssertExecutablePath(t, cmd.Path, NAME)
	if !reflect.DeepEqual(cmd.Args[1:], expectedArgs) {
		t.Fatalf("expected args %v, got %v", expectedArgs, cmd.Args[1:])
	}
}
