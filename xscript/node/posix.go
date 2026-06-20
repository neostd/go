//go:build !windows

package node

import (
	"github.com/neostd/go/exec"
	"github.com/neostd/go/xscript"
)

func init() {
	exec.Register("node", &exec.Executable{
		Name:     "node",
		Variable: xscript.GetVar("NODE"),
		Linux: []string{
			"/usr/bin/node",
			"/usr/local/bin/node",
		},
	})
}
