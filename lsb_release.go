package on_machine

import (
	"os/exec"
	"strings"
)

var LsbRelease *string = nil

func LsbReleaseShell() (*string, error) {
	if LsbRelease != nil {
		return LsbRelease, nil
	}
	path, err := exec.LookPath("lsb_release")
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(path, "-si")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	outStr := strings.ToLower(strings.TrimSpace(string(out)))
	return &outStr, nil
}

func SetLsbReleaseShell() error {
	lsbRelease, err := LsbReleaseShell()
	if err != nil {
		return err
	}
	LsbRelease = lsbRelease
	return nil
}
