//go:build !windows

package nushell

import (
	"github.com/neostd/go/exec"
	"github.com/neostd/go/xscript"
)

func init() {
	exec.Register("nu", &exec.Executable{
		Name:     "nu",
		Variable: xscript.GetVar("NU"),
		Linux: []string{
			"/bin/nu",
			"/usr/bin/nu",
			"/usr/local/bin/nu",
		},
	})
}
