//go:build windows

package zsh

import (
	"github.com/neostd/go/exec"
	"github.com/neostd/go/xscript"
)

func init() {
	exec.Register("zsh", &exec.Executable{
		Name:     "zsh",
		Variable: xscript.GetVar("ZSH"),
		Windows: []string{
			"${SystemRoot}\\msys64\\usr\\bin\\zsh.exe",
			"${SystemRoot}\\..\\msys64\\usr\\bin\\zsh.exe",
			"${SystemRoot}\\..\\cygwin64\\bin\\zsh.exe",
			"${SystemRoot}\\..\\cygwin\\bin\\zsh.exe",
		},
	})
}
