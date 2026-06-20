//go:build windows

package bash

import (
	"path/filepath"
	"strings"

	"github.com/neostd/go/exec"
	"github.com/neostd/go/xscript"
)

func init() {
	exec.Register("bash", &exec.Executable{
		Name:     "bash",
		Variable: xscript.GetVar("BASH"),
		Windows: []string{
			"${ProgramFiles}\\Git\\bin\\bash.exe",
			"${ProgramFiles(x86)}\\Git\\bin\\bash.exe",
			"${ProgramFiles}\\Git\\usr\\bin\\bash.exe",
			"${ProgramFiles(x86)}\\Git\\usr\\bin\\bash.exe",
			"${SystemRoot}\\msys64\\usr\\bin\\bash.exe",
			"${SystemRoot}\\System32\\bash.exe",
		},
	})

	exec.Registry.Get("bash")
}

func resolveScriptFile(script string) string {
	if !filepath.IsAbs(script) {
		file, err := filepath.Abs(script)
		if err != nil {
			script = file
		}
	}

	// determine if bash is the WSL one.
	bash, _ := exec.Find("bash", nil)
	if !strings.Contains(strings.ToLower(bash), "system32") {
		return script
	}

	script = strings.ReplaceAll(script, "\\", "/")
	if script[1] == ':' {
		script = "/mnt/" + strings.ToLower(script[0:1]) + script[2:]
	}

	return script
}
