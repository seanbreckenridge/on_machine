package on_machine

import (
	"errors"
	"strings"
)

// implements a printf-like syntax:
//
// %o - Operating System (using uname)
// %d - Distro (using lsb_release)
// %h - Hostname (name of the computer)
// %a - Arch (detected by golang)
// %O - Golang OS (unmodified golang detected operating system)

func replaceField(pat string) (string, error) {
	switch pat {
	case "%o":
		return GetOS(), nil
	case "%a":
		return GetGolangArch(), nil
	case "%d":
		return GetDistro(), nil
	case "%h":
		return GetHostname(), nil
	case "%O":
		return GetGolangOS(), nil
	default:
		return "", errors.New("Did not match format specifier")
	}
}

func ReplaceFields(pattern string) string {
	// if length 0 or 1, returns pattern
	if len(pattern) <= 1 {
		return pattern
	}
	var sb strings.Builder
	i := 0
	for i < len(pattern) {
		// make sure we don't go out of bounds
		endIndex := min(i+2, len(pattern))
		chunk := pattern[i:endIndex]
		matched, err := replaceField(chunk)
		if err == nil {
			// if matched a format specifier, increment by two
			sb.WriteString(matched)
			i += 2
		} else {
			sb.WriteString(pattern[i : i+1])
			i += 1
		}
	}
	return sb.String()
}
