//go:build !windows

package elixir

import (
	"github.com/neostd/go/exec"
	"github.com/neostd/go/xscript"
)

func init() {
	exec.Register("elixir", &exec.Executable{
		Name:     "elixir",
		Variable: xscript.GetVar("ELIXIR"),
		Linux: []string{
			"/usr/bin/elixir",
			"/usr/local/bin/elixir",
		},
	})
}
