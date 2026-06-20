//go:build !windows

package bun

import (
	"github.com/neostd/go/exec"
	"github.com/neostd/go/xscript"
)

func init() {
	exec.Register("bun", &exec.Executable{
		Name:     "bun",
		Variable: xscript.GetVar("BUN"),
		Linux: []string{
			"${HOME}/.bun/bin/bun",
			"${HOME}/.local/bin/bun",
			"/usr/bin/bun",
			"/usr/local/bin/bun",
		},
	})
}
