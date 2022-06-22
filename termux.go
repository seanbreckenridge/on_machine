package on_machine

import (
	"os"
	"os/exec"
)

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

func OnTermux() bool {
	res, _, _ := Cache.Memoize("on-termux", func() (interface{}, error) {
		return commandExists("termux-setup-storage") && dirExists("/system/app/") && dirExists("/system/priv-app/"), nil
	})
	if castRes, ok := res.(bool); ok {
		return castRes
	} else {
		return false
	}
}
