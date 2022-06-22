package on_machine

import (
	"os"
	"runtime"
	"strings"
)

func GetGolangOS() string {
	return runtime.GOOS
}

func GetGolangArch() string {
	return runtime.GOARCH
}

func GetOS() string {
	os := GetGolangOS()
	uname, _ := UnameSh()
	if uname != nil {
		// ported from neofetch
		switch uname.name {
		case "Darwin":
			return "mac"
		case "SunOS":
			return "solaris"
		case "Haiku":
			return "haiku"
		case "AIX":
			return "aix"
		case "MINIX":
			return "minix"
		case "FreeMiNT":
			return "freemint"
		}
		if uname.name == "Linux" || strings.HasPrefix(uname.name, "GNU") {
			if OnTermux() {
				return "android"
			}
			return "linux"
		}
		if uname.name == "DragonFly" || uname.name == "Bitrig" || strings.HasSuffix(uname.name, "BSD") {
			return "bsd"
		}
		if strings.HasPrefix(uname.name, "CYGWIN") || strings.HasPrefix(uname.name, "MSYS") || strings.HasPrefix(uname.name, "MINGW") {
			return "windows"
		}
	}
	return os
}

func GetDistro() (distro string) {
	uname, _ := UnameSh()
	lsbRelease, _ := LsbReleaseSh()
	distro = "unknown"
	if lsbRelease != nil {
		distro = *lsbRelease
	} else if uname != nil {
		distro = uname.version
	}
	if OnTermux() {
		return "termux"
	}
	return
}

func Hostname() (*string, error) {
	res, err, _ := Cache.Memoize("hostname", func() (interface{}, error) {
		name, err := os.Hostname()
		if err != nil {
			return nil, err
		}
		return name, nil
	})
	if castRes, ok := res.(string); ok {
		return &castRes, nil
	} else {
		return nil, err
	}
}

func GetHostname() string {
	host, err := Hostname()
	if err == nil {
		return *host
	} else {
		return ""
	}
}
