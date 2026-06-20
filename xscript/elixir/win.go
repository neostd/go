//go:build windows

package elixir

import (
	"github.com/neostd/go/exec"
	"github.com/neostd/go/xscript"
)

func init() {
	exec.Register("elixir", &exec.Executable{
		Name:     "elixir",
		Variable: xscript.GetVar("ELIXIR"),
		Windows: []string{
			"${ProgramFiles}\\Elixir\\bin\\elixir.bat",
			"${ProgramFiles(x86)}\\Elixir\\bin\\elixir.bat",
			"${CHOCOLATEYINSTALL}\\bin\\elixir.bat",
		},
	})
}
