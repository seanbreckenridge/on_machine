package on_machine

import (
	"io/ioutil"
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

	res, err, _ := Cache.Memoize("os", func() (interface{}, error) {
		os := GetGolangOS()
		uname, _ := UnameSh()
		if strings.Contains(strings.ToLower(uname.version), "microsoft") {
			return "windows", nil
		}
		if ok, _ := PathExists("/proc/version"); ok {
			contents, err := ioutil.ReadFile("/proc/version")
			if err != nil {
				version := string(contents)
				if strings.Contains(strings.ToLower(version), "microsoft") {
					return "windows", nil
				}
			}
		}
		if uname != nil {
			// ported from neofetch
			switch uname.name {
			case "Darwin":
				return "mac", nil
			case "SunOS":
				return "solaris", nil
			case "Haiku":
				return "haiku", nil
			case "AIX":
				return "aix", nil
			case "MINIX":
				return "minix", nil
			case "FreeMiNT":
				return "freemint", nil
			}
			if uname.name == "Linux" || strings.HasPrefix(uname.name, "GNU") {
				if OnTermux() {
					return "android", nil
				}
				return "linux", nil
			}
			if uname.name == "DragonFly" || uname.name == "Bitrig" || strings.HasSuffix(uname.name, "BSD") {
				return "bsd", nil
			}
			if strings.HasPrefix(uname.name, "CYGWIN") || strings.HasPrefix(uname.name, "MSYS") || strings.HasPrefix(uname.name, "MINGW") {
				return "windows", nil
			}
		}
		return strings.ToLower(os), nil
	})
	if castRes, ok := res.(string); ok {
		return castRes
	} else {
		panic(err)
	}
}

func GetDistro() string {
	if OnTermux() {
		return "termux"
	}
	if GetOS() == "windows" {
		// I can't see a case where this is detected as windows but not in WSL?
		return "wsl"
	}
	lsbRelease, _ := LsbReleaseSh()
	if lsbRelease != nil {
		return *lsbRelease
	}
	uname, _ := UnameSh()
	if uname != nil {
		return uname.version
	}
	return "unknown"
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
		return "unknown"
	}
}
