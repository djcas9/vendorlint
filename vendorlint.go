package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mephux/vendorlint/vendorlint"
)

const (
	// Name of the application
	Name = "vendorlint"

	// Version of the application
	Version = "0.1.0"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr,
			fmt.Sprintf("Usage: %s (%s) [-t] package [package ...]", Name, Version),
		)

		flag.PrintDefaults()
	}

	tests := flag.Bool("t", false, "include test dependencies")
	version := flag.Bool("v", false, fmt.Sprintf("%s version number", Name))

	flag.Parse()

	if *version {
		fmt.Fprintln(os.Stdout, Version)
		os.Exit(0)
	}

	if len(flag.Args()) <= 0 {
		flag.Usage()
		os.Exit(0)
	}

	config := vendorlint.NewConfig()
	config.Tests = *tests
	config.Packages = flag.Args()

	if cwd, err := os.Getwd(); err == nil {
		config.WorkingDirectory = cwd
	} else {
		// error!
	}

	if linter, err := vendorlint.NewLinter(config); err == nil {
		linter.Report()
	} else {
		// error!
	}
}
