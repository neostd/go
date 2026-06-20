//go:build !windows

package lua

import (
	"github.com/neostd/go/exec"
	"github.com/neostd/go/xscript"
)

func init() {
	exec.Register("lua", &exec.Executable{
		Name:     "lua",
		Variable: xscript.GetVar("LUA"),
		Linux: []string{
			"/usr/bin/lua",
			"/usr/local/bin/lua",
		},
	})
}
