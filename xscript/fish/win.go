//go:build windows

package fish

import (
	"github.com/neostd/go/exec"
	"github.com/neostd/go/xscript"
)

func init() {
	exec.Register("fish", &exec.Executable{
		Name:     "fish",
		Variable: xscript.GetVar("FISH"),
		Windows: []string{
			"${SystemRoot}\\msys64\\usr\\bin\\fish.exe",
			"${SystemRoot}\\..\\msys64\\usr\\bin\\fish.exe",
			"${SystemRoot}\\..\\cygwin64\\bin\\fish.exe",
			"${SystemRoot}\\..\\cygwin\\bin\\fish.exe",
		},
	})
}
