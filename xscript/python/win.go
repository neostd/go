//go:build windows

package python

import (
	"github.com/neostd/go/exec"
	"github.com/neostd/go/xscript"
)

func init() {
	exec.Register("python", &exec.Executable{
		Name:     "python",
		Variable: xscript.GetVar("PYTHON"),
		Windows: []string{
			"${ProgramFiles}\\Python\\Python.exe",
			"${ProgramFiles(x86)}\\Python\\Python.exe",
		},
	})
}
