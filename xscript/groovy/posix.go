//go:build !windows

package groovy

import (
	"github.com/neostd/go/exec"
	"github.com/neostd/go/xscript"
)

func init() {
	exec.Register("groovy", &exec.Executable{
		Name:     "groovy",
		Variable: xscript.GetVar("GROOVY"),
		Linux: []string{
			"/usr/bin/groovy",
			"/usr/local/bin/groovy",
		},
	})
}
