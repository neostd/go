//go:build windows

package dotnet

import (
	"github.com/neostd/go/exec"
	"github.com/neostd/go/xscript"
)

func init() {
	exec.Register("dotnet", &exec.Executable{
		Name:     "dotnet",
		Variable: xscript.GetVar("DOTNET"),
		Windows: []string{
			"${LOCALAPPDATA}\\Microsoft\\dotnet\\dotnet.exe",
			"${USERPROFILE}\\.dotnet\\dotnet.exe",
			"${USERPROFILE}\\.local\\bin\\dotnet.exe",
			"${ProgramFiles}\\dotnet\\dotnet.exe",
		},
	})
}
