//go:build !windows

package dotnet

import (
	"github.com/neostd/go/exec"
	"github.com/neostd/go/xscript"
)

func init() {
	exec.Register("dotnet", &exec.Executable{
		Name:     "dotnet",
		Variable: xscript.GetVar("DOTNET"),
		Linux: []string{
			"$HOME/.dotnet/dotnet",
			"$HOME/.local/share/dotnet/dotnet",
			"/usr/local/bin/dotnet",
			"/usr/bin/dotnet",
			"/bin/dotnet",
		},
	})
}
