//go:build !windows

package bash

import (
	"github.com/neostd/go/exec"
	"github.com/neostd/go/xscript"
)

func init() {
	exec.Register("bash", &exec.Executable{
		Name:     "bash",
		Variable: xscript.GetVar("BASH"),
		Linux: []string{
			"/bin/bash",
			"/usr/bin/bash",
		},
	})
}

func resolveScriptFile(script string) string {
	return script
}
