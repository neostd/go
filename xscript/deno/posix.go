//go:build !windows

package deno

import (
	"github.com/neostd/go/exec"
	"github.com/neostd/go/xscript"
)

func init() {
	exec.Register("deno", &exec.Executable{
		Name:     "deno",
		Variable: xscript.GetVar("DENO"),
		Linux: []string{
			"/usr/bin/deno",
			"/usr/local/bin/deno",
		},
	})
}
