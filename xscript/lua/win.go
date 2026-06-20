//go:build windows

package lua

import (
	"github.com/neostd/go/exec"
	"github.com/neostd/go/xscript"
)

func init() {
	exec.Register("lua", &exec.Executable{
		Name:     "lua",
		Variable: xscript.GetVar("LUA"),
		Windows: []string{
			"${ProgramFiles}\\Lua\\lua.exe",
			"${ProgramFiles(x86)}\\Lua\\lua.exe",
			"${CHOCOLATEYINSTALL}\\bin\\lua.exe",
		},
	})
}
