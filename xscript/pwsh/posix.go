//go:build !windows

package pwsh

import (
	"github.com/neostd/go/exec"
	"github.com/neostd/go/xscript"
)

func init() {
	exec.Register("pwsh", &exec.Executable{
		Name:     "pwsh",
		Variable: xscript.GetVar("PWSH"),
		Linux: []string{
			"/bin/pwsh",
			"/usr/bin/pwsh",
			"/opt/microsoft/powershell/7/pwsh",
			"/opt/microsoft/powershell/6/pwsh",
		},
	})
}
