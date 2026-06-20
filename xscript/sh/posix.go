//go:build !windows

package sh

import (
	"github.com/neostd/go/exec"
	"github.com/neostd/go/xscript"
)

func init() {
	exec.Register("sh", &exec.Executable{
		Name:     "sh",
		Variable: xscript.GetVar("SH"),
		Linux: []string{
			"/bin/sh",
			"/usr/bin/sh",
		},
	})
}
