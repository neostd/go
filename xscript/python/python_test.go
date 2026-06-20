package python

import (
	"reflect"
	"testing"

	"github.com/neostd/go/xscript/internal/testutil"
)

func init() {
	testutil.SetupMise("python", "python", "python3")
}

func TestNew(t *testing.T) {
	testutil.RequireExecutable(t, NAME)

	cmd := New("arg1", "arg2")
	expectedArgs := []string{"arg1", "arg2"}
	testutil.AssertExecutablePath(t, cmd.Path, NAME, "python3")
	if !reflect.DeepEqual(cmd.Args[1:], expectedArgs) {
		t.Fatalf("expected args %v, got %v", expectedArgs, cmd.Args[1:])
	}
}

func TestFile(t *testing.T) {
	testutil.RequireExecutable(t, NAME, "python3")

	cmd := File("test.py", "arg1")
	expectedArgs := []string{"test.py", "arg1"}
	testutil.AssertExecutablePath(t, cmd.Path, NAME, "python3")
	if !reflect.DeepEqual(cmd.Args[1:], expectedArgs) {
		t.Fatalf("expected args %v, got %v", expectedArgs, cmd.Args[1:])
	}
}

func TestInline(t *testing.T) {
	testutil.RequireExecutable(t, NAME, "python3")

	cmd := Inline("print('hello')", "arg1")
	expectedArgs := []string{"-c", "print('hello')", "arg1"}
	testutil.AssertExecutablePath(t, cmd.Path, NAME, "python3")
	if !reflect.DeepEqual(cmd.Args[1:], expectedArgs) {
		t.Fatalf("expected args %v, got %v", expectedArgs, cmd.Args[1:])
	}
}

func TestScriptFile(t *testing.T) {
	testutil.RequireExecutable(t, NAME, "python3")

	cmd := Script("test.py", "arg1")
	expectedArgs := []string{"test.py", "arg1"}
	testutil.AssertExecutablePath(t, cmd.Path, NAME, "python3")
	if !reflect.DeepEqual(cmd.Args[1:], expectedArgs) {
		t.Fatalf("expected args %v, got %v", expectedArgs, cmd.Args[1:])
	}
}

func TestScriptInline(t *testing.T) {
	testutil.RequireExecutable(t, NAME, "python3")

	cmd := Script("print('hello')", "arg1")
	expectedArgs := []string{"-c", "print('hello')", "arg1"}
	testutil.AssertExecutablePath(t, cmd.Path, NAME, "python3")
	if !reflect.DeepEqual(cmd.Args[1:], expectedArgs) {
		t.Fatalf("expected args %v, got %v", expectedArgs, cmd.Args[1:])
	}
}
