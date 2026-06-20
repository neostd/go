//go:build windows

package sh

import (
	"github.com/neostd/go/exec"
	"github.com/neostd/go/xscript"
)

func init() {
	exec.Register("sh", &exec.Executable{
		Name:     "sh",
		Variable: xscript.GetVar("SH"),
		Windows: []string{
			"${ProgramFiles}\\Git\\bin\\sh.exe",
			"${ProgramFiles(x86)}\\Git\\bin\\sh.exe",
		},
	})
}
