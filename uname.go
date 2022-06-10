package on_machine

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// For possible GolangArch/GolangOS possibilities, run 'go tool dist list'

type UnameOutput struct {
	name    string
	version string
	machine string
}

var Uname *UnameOutput = nil

func UnameShell() (*UnameOutput, error) {
	if Uname != nil {
		return Uname, nil
	}
	// if this has already been called, don't do so again
	path, err := exec.LookPath("uname")
	if err != nil {
		return nil, err
	}
	// call and grab name/version/machine from uname
	// hopefully this is portable - seems to be what neofetch does
	cmd := exec.Command(path, "-srm")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	// split into the 3 fields
	output := string(out)
	parts := strings.Fields(output)
	if len(parts) != 3 {
		return nil, errors.New(fmt.Sprintf("uname -srm didn't successfully split into 3 items: %s\n", output))
	}
	uname := UnameOutput{name: parts[0], version: parts[1], machine: parts[2]}
	return &uname, nil
}

func SetUnameShell() error {
	uname, err := UnameShell()
	if err != nil {
		return err
	}
	Uname = uname
	return nil
}
