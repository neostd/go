//go:build windows

package deno

import (
	"github.com/neostd/go/exec"
	"github.com/neostd/go/xscript"
)

func init() {
	exec.Register("deno", &exec.Executable{
		Name:     "deno",
		Variable: xscript.GetVar("DENO"),
		Windows: []string{
			"${USERPROFILE}\\.deno\\bin\\deno.exe",
			"${LOCALAPPDATA}\\Programs\\deno\\deno.exe",
			"${CHOCOLATEYINSTALL}\\bin\\deno.exe",
			"${USERPROFILE}\\.local\\share\\bin\\deno.exe",
			"${ProgramFiles}\\deno\\bin\\deno.exe",
		},
	})
}
