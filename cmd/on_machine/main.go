package main

import (
	"flag"
	"fmt"
	"github.com/seanbreckenridge/on_machine"
	"os"
)

type OnMachineConfig struct {
	pattern string
}

func parseFlags() (*OnMachineConfig, error) {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `usage: on_machine [-h] [PATTERN]

Tool to determine which operating system/machine you're on.

PATTERN is a printf-styled format string, supporting the following sequences:

%o - Operating System (using uname)
%d - Distro (using lsb_release)
%h - Hostname (name of the computer)
%a - Arch (detected by golang)
%O - Golang OS (unmodified golang detected operating system)

By default, this uses '%o_%d_%h'`)
		flag.PrintDefaults()
	}
	flag.Parse()
	var pattern string
	switch flag.NArg() {
	case 1:
		pattern = flag.Arg(0)
	default:
		pattern = "%o_%d_%h"
	}
	return &OnMachineConfig{
		pattern: pattern,
	}, nil
}

func run() error {
	conf, err := parseFlags()
	if err != nil {
		return err
	}
	res := on_machine.ReplaceFields(conf.pattern)
	fmt.Println(res)
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
