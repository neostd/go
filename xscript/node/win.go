//go:build windows

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
			"${USERPROFILE}\\.node\\bin\\node.exe",
			"${USERPROFILE}\\.local\\bin\\node.exe",
			"${ProgramFiles}\\nodejs\\node.exe",
		},
	})
}
