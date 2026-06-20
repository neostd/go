//go:build windows

package nushell

import (
	"github.com/neostd/go/exec"
	"github.com/neostd/go/xscript"
)

func init() {
	exec.Register("nu", &exec.Executable{
		Name:     "nu",
		Variable: xscript.GetVar("NU"),
		Windows: []string{
			"${ProgramFiles}\\nu\\bin\\nu.exe",
			"${LocalAppData}\\Programs\\nu\\bin\\nu.exe",
		},
	})
}
