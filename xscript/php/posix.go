//go:build !windows

package php

import (
	"github.com/neostd/go/exec"
	"github.com/neostd/go/xscript"
)

func init() {
	exec.Register("php", &exec.Executable{
		Name:     "php",
		Variable: xscript.GetVar("PHP"),
		Linux: []string{
			"/usr/bin/php",
			"/usr/local/bin/php",
		},
	})
}
