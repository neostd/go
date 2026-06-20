//go:build !windows

package python

import (
	"github.com/neostd/go/exec"
	"github.com/neostd/go/xscript"
)

func init() {
	exec.Register("python", &exec.Executable{
		Name:     "python",
		Variable: xscript.GetVar("PYTHON"),
		Linux: []string{
			"/usr/bin/python",
			"/usr/bin/python3",
			"/usr/local/bin/python",
			"/usr/local/bin/python3",
		},
	})
}
