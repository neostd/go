package lua

import (
	"reflect"
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

	cmd := File("test.lua", "arg1")
	expectedArgs := []string{"test.lua", "arg1"}
	testutil.AssertExecutablePath(t, cmd.Path, NAME)
	if !reflect.DeepEqual(cmd.Args[1:], expectedArgs) {
		t.Fatalf("expected args %v, got %v", expectedArgs, cmd.Args[1:])
	}
}

func TestInline(t *testing.T) {
	testutil.RequireExecutable(t, NAME)

	cmd := Inline("print('hello')", "arg1")
	expectedArgs := []string{"-e", "print('hello')", "arg1"}
	testutil.AssertExecutablePath(t, cmd.Path, NAME)
	if !reflect.DeepEqual(cmd.Args[1:], expectedArgs) {
		t.Fatalf("expected args %v, got %v", expectedArgs, cmd.Args[1:])
	}
}

func TestScriptFile(t *testing.T) {
	testutil.RequireExecutable(t, NAME)

	cmd := Script("test.lua", "arg1")
	expectedArgs := []string{"test.lua", "arg1"}
	testutil.AssertExecutablePath(t, cmd.Path, NAME)
	if !reflect.DeepEqual(cmd.Args[1:], expectedArgs) {
		t.Fatalf("expected args %v, got %v", expectedArgs, cmd.Args[1:])
	}
}

func TestScriptInline(t *testing.T) {
	testutil.RequireExecutable(t, NAME)

	cmd := Script("print('hello')", "arg1")
	expectedArgs := []string{"-e", "print('hello')", "arg1"}
	testutil.AssertExecutablePath(t, cmd.Path, NAME)
	if !reflect.DeepEqual(cmd.Args[1:], expectedArgs) {
		t.Fatalf("expected args %v, got %v", expectedArgs, cmd.Args[1:])
	}
}
