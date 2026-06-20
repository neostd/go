package powershell

import (
	"reflect"
	"testing"

	"github.com/neostd/go/xscript/internal/testutil"
)

func init() {
	testutil.SetupMise(NAME, "pwsh", NAME)
}

func TestNew(t *testing.T) {
	testutil.RequireExecutable(t, NAME, "pwsh")

	cmd := New("arg1", "arg2")
	expectedArgs := []string{"arg1", "arg2"}
	testutil.AssertExecutablePath(t, cmd.Path, NAME, "pwsh")
	if !reflect.DeepEqual(cmd.Args[1:], expectedArgs) {
		t.Fatalf("expected args %v, got %v", expectedArgs, cmd.Args[1:])
	}
}

func TestFile(t *testing.T) {
	testutil.RequireExecutable(t, NAME, "pwsh")

	cmd := File("test.ps1", "arg1")
	expectedArgs := []string{"-NoLogo", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", "test.ps1", "arg1"}
	testutil.AssertExecutablePath(t, cmd.Path, NAME, "pwsh")
	if !reflect.DeepEqual(cmd.Args[1:], expectedArgs) {
		t.Fatalf("expected args %v, got %v", expectedArgs, cmd.Args[1:])
	}
}

func TestInline(t *testing.T) {
	testutil.RequireExecutable(t, NAME, "pwsh")

	cmd := Inline("Write-Output hello", "arg1")
	expectedArgs := []string{"-NoLogo", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", "Write-Output hello", "arg1"}
	testutil.AssertExecutablePath(t, cmd.Path, NAME, "pwsh")
	if !reflect.DeepEqual(cmd.Args[1:], expectedArgs) {
		t.Fatalf("expected args %v, got %v", expectedArgs, cmd.Args[1:])
	}
}

func TestScriptFile(t *testing.T) {
	testutil.RequireExecutable(t, NAME, "pwsh")

	cmd := Script("test.ps1", "arg1")
	expectedArgs := []string{"-NoLogo", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", "test.ps1", "arg1"}
	testutil.AssertExecutablePath(t, cmd.Path, NAME, "pwsh")
	if !reflect.DeepEqual(cmd.Args[1:], expectedArgs) {
		t.Fatalf("expected args %v, got %v", expectedArgs, cmd.Args[1:])
	}
}

func TestScriptInline(t *testing.T) {
	testutil.RequireExecutable(t, NAME, "pwsh")

	cmd := Script("Write-Output hello", "arg1")
	expectedArgs := []string{"-NoLogo", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", "Write-Output hello", "arg1"}
	testutil.AssertExecutablePath(t, cmd.Path, NAME, "pwsh")
	if !reflect.DeepEqual(cmd.Args[1:], expectedArgs) {
		t.Fatalf("expected args %v, got %v", expectedArgs, cmd.Args[1:])
	}
}
