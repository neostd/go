// Package testutil contains shared helpers used by xscript tests.
package testutil

import (
	"bytes"
	"os"
	stdexec "os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	neosexec "github.com/neostd/go/exec"
)

var pathMu sync.Mutex

// RequireExecutable skips the test unless one of the provided executable names
// can be resolved on PATH.
func RequireExecutable(t *testing.T, names ...string) string {
	t.Helper()

	for _, name := range names {
		path, err := stdexec.LookPath(name)
		if err == nil {
			return path
		}
	}

	t.Skipf("skipping test: executable not found on PATH: %s", strings.Join(names, ", "))
	return ""
}

// SetupMise installs a runtime with mise and prepends its bin directory to
// PATH when USE_MISE_FOR_TESTING=true. USE_MISE_FOR_TESTS is also accepted
// for backward compatibility.
func SetupMise(tool string, bins ...string) {
	SetupMiseWithDeps(tool, nil, bins...)
}

// SetupMiseWithDeps installs any prerequisite tools before installing the
// requested runtime with mise.
func SetupMiseWithDeps(tool string, deps []string, bins ...string) {
	useMise := strings.ToLower(os.Getenv("USE_MISE_FOR_TESTING")) == "true" ||
		strings.ToLower(os.Getenv("USE_MISE_FOR_TESTS")) == "true"
	if !useMise {
		return
	}

	misePath, ok := neosexec.Which("mise")
	if !ok || misePath == "" {
		return
	}

	for _, dep := range deps {
		if !installMiseTool(misePath, dep+"@latest") {
			return
		}
	}

	toolSpec := tool + "@latest"
	if !installMiseTool(misePath, toolSpec) {
		return
	}

	for _, bin := range bins {
		which := stdexec.Command(misePath, "which", bin, "--tool", toolSpec)
		output, err := which.Output()
		if err != nil {
			continue
		}

		path := strings.TrimSpace(string(output))
		if path == "" {
			continue
		}

		prependPath(filepath.Dir(path))
		return
	}
}

func installMiseTool(misePath, toolSpec string) bool {
	install := stdexec.Command(misePath, "install", "-y", toolSpec)
	install.Stdout = &bytes.Buffer{}
	install.Stderr = &bytes.Buffer{}
	return install.Run() == nil
}

// prependPath adds a directory to the front of PATH once per process.
func prependPath(dir string) {
	if dir == "" {
		return
	}

	pathMu.Lock()
	defer pathMu.Unlock()

	parts := filepath.SplitList(os.Getenv("PATH"))
	for _, part := range parts {
		if strings.EqualFold(part, dir) {
			return
		}
	}

	if len(parts) == 0 || parts[0] == "" {
		_ = os.Setenv("PATH", dir)
		return
	}

	_ = os.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

// AssertExecutablePath verifies that the resolved executable matches one of the
// expected runtime names, accounting for common Windows suffixes.
func AssertExecutablePath(t *testing.T, path string, names ...string) {
	t.Helper()

	base := strings.ToLower(filepath.Base(path))
	trimmed := strings.TrimSuffix(strings.TrimSuffix(strings.TrimSuffix(base, ".exe"), ".bat"), ".cmd")
	for _, name := range names {
		lower := strings.ToLower(name)
		if base == lower || trimmed == lower {
			return
		}
	}

	t.Fatalf("expected executable path %q to match one of %v", path, names)
}
