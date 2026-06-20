//go:build windows
// +build windows

package xscript

func GetVar(name string) string {
	varName := ""
	if prefix != "" {
		varName = prefix + "_" + name
	} else {
		varName = name
	}

	if win != "" {
		varName = varName + "_" + win
	}

	if suffix != "" {
		varName = varName + "_" + suffix
	}

	return varName
}
