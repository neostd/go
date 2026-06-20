//go:build windows

package golang

import (
	"github.com/neostd/go/exec"
	"github.com/neostd/go/xscript"
)

func init() {
	exec.Register("go", &exec.Executable{
		Name:     "go",
		Variable: xscript.GetVar("GO"),
		Windows: []string{
			"${ProgramFiles}\\Go\\bin\\go.exe",
			"${ChocolateyInstall}\\lib\\go\\tools\\go\\bin\\go.exe",
		},
	})
}
