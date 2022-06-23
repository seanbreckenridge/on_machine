package main

import (
	"flag"
	"fmt"
	"github.com/seanbreckenridge/on_machine"
	"os"
)

const (
	PRINT       = 1
	MATCH_PATHS = 2
)

type Command = int

type OnMachineConfig struct {
	pattern      string
	matchBaseDir string
	command      Command
}

func parseFlags() (*OnMachineConfig, error) {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `usage: on_machine [-h] [-cmd <print|match>] [OPTIONS] [PATTERN]

Tool to determine which operating system/machine you're on.

Commands:
print [default]: prints the resulting pattern after interpolating the pattern
match: does directory/path matching based on the pattern, changes the default pattern to '%o/%d/%h'

print
---
PATTERN is a printf-styled format string, supporting the following sequences:

%o - Operating System (using uname)
%d - Distro (using lsb_release)
%h - Hostname (name of the computer)
%a - Arch (detected by golang)
%O - Golang OS (unmodified golang detected operating system)

By default, this uses '%o_%d_%h'

match
---
Directory/path matching, Uses the pattern to match directory structures.
Can provide the base path to use with -base, that replaces '/' with
OS-specific path separator in the pattern. For more information, see the docs:
https://github.com/seanbreckenridge/on_machine

Options:
`)
		flag.PrintDefaults()
	}
	cmd := flag.String("cmd", "print", "on_machine command to run")
	base := flag.String("base", "", "Base directory to use to match paths")
	flag.Parse()
	var pattern string
	// parse command
	var command Command
	switch *cmd {
	case "print":
		command = PRINT
	case "match":
		command = MATCH_PATHS
	default:
		fmt.Printf("Unknown command '%s'. Provide either 'print' or 'match'\n", *cmd)
		os.Exit(1)
	}
	// set pattern
	switch flag.NArg() {
	case 1:
		pattern = flag.Arg(0)
	default:
		// set default pattern
		switch command {
		case PRINT:
			pattern = "%o_%d_%h"
		case MATCH_PATHS:
			pattern = "%o/%d/%h"
		}
	}
	// match based parsing
	var matchBase string
	if command == MATCH_PATHS {
		matchBase = string(*base)
		if matchBase != "" {
			if !on_machine.DirExists(matchBase) {
				fmt.Fprintf(os.Stderr, "Directory doesnt exist: '%s'\n", matchBase)
				os.Exit(1)
			}
		}
	}
	return &OnMachineConfig{
		pattern:      pattern,
		command:      command,
		matchBaseDir: matchBase,
	}, nil
}

func run() error {
	conf, err := parseFlags()
	if err != nil {
		return err
	}
	switch conf.command {
	case PRINT:
		res := on_machine.ReplaceFields(conf.pattern)
		fmt.Println(res)
	case MATCH_PATHS:
		matched, _ := on_machine.MatchPaths(conf.pattern, conf.matchBaseDir)
		fmt.Printf("%+v\n", matched)
	}
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
