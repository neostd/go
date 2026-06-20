//go:build !windows

package zsh

import (
	"github.com/neostd/go/exec"
	"github.com/neostd/go/xscript"
)

func init() {
	exec.Register("zsh", &exec.Executable{
		Name:     "zsh",
		Variable: xscript.GetVar("ZSH"),
		Linux: []string{
			"/bin/zsh",
			"/usr/bin/zsh",
			"/usr/local/bin/zsh",
		},
	})
}
