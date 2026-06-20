//go:build !windows

package powershell

import (
	"github.com/neostd/go/exec"
	"github.com/neostd/go/xscript"
)

func init() {
	exec.Register("powershell", &exec.Executable{
		Name:     "powershell",
		Variable: xscript.GetVar("POWERSHELL"),
		Linux: []string{
			"/bin/pwsh",
			"/usr/bin/pwsh",
		},
	})
}
