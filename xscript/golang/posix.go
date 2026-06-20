//go:build !windows

package golang

import (
	"github.com/neostd/go/exec"
	"github.com/neostd/go/xscript"
)

func init() {
	exec.Register("go", &exec.Executable{
		Name:     "go",
		Variable: xscript.GetVar("GO"),
		Linux: []string{
			"${HOME}/.local/shared/go/bin/go",
			"/usr/local/go/bin/go",
			"/usr/local/bin/go",
			"/usr/bin/go",
		},
	})
}
