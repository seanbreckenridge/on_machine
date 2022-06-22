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

func UnameSh() (*UnameOutput, error) {
	res, err, _ := Cache.Memoize("uname", func() (interface{}, error) {
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
	})
	// type assert
	if castRes, ok := res.(*UnameOutput); ok {
		return castRes, nil
	} else {
		return nil, err
	}
}
