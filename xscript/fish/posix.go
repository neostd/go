//go:build !windows

package fish

import (
	"github.com/neostd/go/exec"
	"github.com/neostd/go/xscript"
)

func init() {
	exec.Register("fish", &exec.Executable{
		Name:     "fish",
		Variable: xscript.GetVar("FISH"),
		Linux: []string{
			"/bin/fish",
			"/usr/bin/fish",
			"/usr/local/bin/fish",
		},
	})
}
