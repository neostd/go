//go:build windows

package php

import (
	"github.com/neostd/go/exec"
	"github.com/neostd/go/xscript"
)

func init() {
	exec.Register("php", &exec.Executable{
		Name:     "php",
		Variable: xscript.GetVar("PHP"),
		Windows: []string{
			"${ProgramFiles}\\PHP\\php.exe",
			"${ProgramFiles(x86)}\\PHP\\php.exe",
			"${CHOCOLATEYINSTALL}\\bin\\php.exe",
		},
	})
}
