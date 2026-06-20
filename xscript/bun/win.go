//go:build windows

package bun

import (
	"github.com/neostd/go/exec"
	"github.com/neostd/go/xscript"
)

func init() {

	exec.Register("bun", &exec.Executable{
		Name:     "bun",
		Variable: xscript.GetVar("BUN"),
		Windows: []string{
			"${USERPROFILE}\\.bun\\bin\\bun.exe",
			"${LOCALAPPDATA}\\Programs\\bin\\bun.exe",
			"${LOCALAPPDATA}\\Microsoft\\WinGet\\Links\\bin.exe",
			"${CHOCOLATEYINSTALL}\\bin\\bun.exe",
			"${USERPROFILE}\\.local\\share\\bin\\bun.exe",
			"${ProgramFiles}\\bun\\bin\\bun.exe",
			"${ProgramFiles(x86)}\\bun\\bin\\bun.exe",
		},
	})
}
