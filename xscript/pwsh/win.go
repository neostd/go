//go:build windows

package pwsh

import (
	"github.com/neostd/go/exec"
	"github.com/neostd/go/xscript"
)

func init() {
	exec.Register("pwsh", &exec.Executable{
		Name:     "pwsh",
		Variable: xscript.GetVar("PWSH"),
		Windows: []string{
			"${ProgramFiles}\\PowerShell\\7\\pwsh.exe",
			"%ProgramFiles(x86)%\\PowerShell\\7\\pwsh.exe",
			"${SystemRoot}\\System32\\WindowsPowerShell\\v1.0\\powershell.exe",
		},
	})
}
