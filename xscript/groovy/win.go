//go:build windows

package groovy

import (
	"github.com/neostd/go/exec"
	"github.com/neostd/go/xscript"
)

func init() {
	exec.Register("groovy", &exec.Executable{
		Name:     "groovy",
		Variable: xscript.GetVar("GROOVY"),
		Windows: []string{
			"${ProgramFiles}\\Groovy\\bin\\groovy.exe",
			"${ProgramFiles(x86)}\\Groovy\\bin\\groovy.exe",
			"${CHOCOLATEYINSTALL}\\bin\\groovy.exe",
		},
	})
}
