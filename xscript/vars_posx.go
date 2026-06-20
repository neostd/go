//go:build !windows
// +build !windows

package xscript

// SetWindows is a no-op on non-Windows platforms and exists to keep the API consistent.
func SetWindows(string) {
}

func GetVar(name string) string {
	varName := ""
	if prefix != "" {
		varName = prefix + "_" + name
	} else {
		varName = name
	}

	if suffix != "" {
		varName = varName + "_" + suffix
	}

	return varName
}
