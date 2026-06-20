//go:build windows

package ruby

import (
	"github.com/neostd/go/exec"
	"github.com/neostd/go/xscript"
)

func init() {
	exec.Register("ruby", &exec.Executable{
		Name:     "ruby",
		Variable: xscript.GetVar("RUBY"),
		Windows: []string{
			"${ProgramFiles}\\Ruby\\bin\\ruby.exe",
			"${ProgramFiles(x86)}\\Ruby\\bin\\ruby.exe",
		},
	})
}
