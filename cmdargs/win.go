//go:build windows
// +build windows

package cmdargs

import "strings"

func containsSpecialChars(s string) bool {
	if len(s) == 0 {
		return false
	}

	for _, c := range s {
		if c == ' ' || c == '"' {
			return true
		}
	}
	return false
}

func appendCliArg(sb *strings.Builder, s string) *strings.Builder {
	if len(s) == 0 {
		return sb
	}

	// basic logic from http://stackoverflow.com/questions/5510343/escape-command-line-arguments-in-c-sharp.
	// and https://blogs.msdn.microsoft.com/twistylittlepassagesallalike/2011/04/23/everyone-quotes-command-line-arguments-the-wrong-way/


	if !containsSpecialChars(s) {
		sb.WriteString(s)
		return sb
	}

	backslashCount := 0
	sb.WriteRune('"')
	for _, c := range s {
		switch c {
		case '\\':
			backslashCount++

		case '"':
			times := (2 * backslashCount) + 1
			backslashCount = 0
			if times > 0 {
				for i := 0; i < times; i++ {
					sb.WriteRune('\\')
				}
			}
			sb.WriteRune('"')
		default:
			if backslashCount > 0 {
				for i := 0; i < backslashCount; i++ {
					sb.WriteRune('\\')
				}
			}
			backslashCount = 0
			sb.WriteRune(c)
		}

	}

	if backslashCount > 0 {
		for i := 0; i < backslashCount; i++ {
			sb.WriteRune('\\')
		}
	}

	sb.WriteRune('"')

	return sb
}
