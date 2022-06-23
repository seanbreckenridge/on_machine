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
		fmt.Fprintln(os.Stderr, `usage: on_machine [-h] [PATTERN]

Tool to determine which operating system/machine you're on.

PATTERN is a printf-styled format string, supporting the following sequences:

%o - Operating System (using uname)
%d - Distro (using lsb_release)
%h - Hostname (name of the computer)
%a - Arch (detected by golang)
%O - Golang OS (unmodified golang detected operating system)

By default, this uses '%o_%d_%h'

-match-paths changes the default pattern to '%o/%d', and uses that pattern to
match directory structures. It expects the directory to use as the 'base'
as the first path, and replaces '/' with the directory separator in the pattern
'
`)
		flag.PrintDefaults()
	}
	matchPaths := flag.String("match-paths", "", "Base directory to use to match paths")
	flag.Parse()
	var pattern string
	command := PRINT
	matchBase := *matchPaths
	if matchBase != "" {
		command = MATCH_PATHS
		pattern = "%o/%d"
		matchBase = *matchPaths
		if !on_machine.DirExists(matchBase) {
			fmt.Fprintf(os.Stderr, "Directory doesnt exist: '%s'\n", matchBase)
			os.Exit(1)
		}
	}
	switch flag.NArg() {
	case 1:
		pattern = flag.Arg(0)
	default:
		pattern = "%o_%d_%h"
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
