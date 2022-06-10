package on_machine

import (
	"os"
	"os/exec"
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
	SetUnameShell()
	os := GetGolangOS()
	if Uname != nil {
		// ported from neofetch
		switch Uname.name {
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
		if Uname.name == "Linux" || strings.HasPrefix(Uname.name, "GNU") {
			if onTermux() {
				return "android"
			}
			return "linux"
		}
		if Uname.name == "DragonFly" || Uname.name == "Bitrig" || strings.HasSuffix(Uname.name, "BSD") {
			return "bsd"
		}
		if strings.HasPrefix(Uname.name, "CYGWIN") || strings.HasPrefix(Uname.name, "MSYS") || strings.HasPrefix(Uname.name, "MINGW") {
			return "windows"
		}
	}
	return os
}

// true if directory exists
func dirExists(path string) bool {
	stat, err := os.Stat(path)
	return err == nil && stat.IsDir()
}

// whether or not the user has a command on their $PATH
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

var OnTermux *bool = nil

func onTermux() bool {
	if OnTermux != nil {
		return *OnTermux
	}
	onTermux := commandExists("termux-setup-storage") && dirExists("/system/app/") && dirExists("/system/priv-app/")
	OnTermux = &onTermux
	return onTermux
}

func GetDistro() (distro string) {
	SetUnameShell()
	SetLsbReleaseShell()
	distro = "unknown"
	if LsbRelease != nil {
		distro = *LsbRelease
	} else if Uname != nil {
		distro = Uname.version
	}
	if onTermux() {
		return "termux"
	}
	return
}

var Hostname *string = nil

func SetHostname() (*string, error) {
	if Hostname != nil {
		return Hostname, nil
	}
	name, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	Hostname = &name
	return Hostname, nil
}

func GetHostname() string {
	hostname, err := SetHostname()
	if err != nil {
		// use localhost as 'unknown value'
		return "localhost"
	} else {
		return *hostname
	}
}
