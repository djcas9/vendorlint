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
	Version = "0.1.2"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr,
			fmt.Sprintf("Usage: %s (%s) [OPTIONS]", Name, Version),
		)

		flag.PrintDefaults()
	}

	tests := flag.Bool("t", false, "include test dependencies")
	missing := flag.Bool("m", false, "report missing dependencies")
	all := flag.Bool("a", false, "report all dependencies")
	paths := flag.Bool("p", false, "only output vendor paths")

	version := flag.Bool("v", false, fmt.Sprintf("%s version number", Name))

	flag.Parse()

	if *version {
		fmt.Fprintln(os.Stdout, Version)
		os.Exit(0)
	}

	config := vendorlint.NewConfig()
	config.Missing = *missing
	config.All = *all
	config.Tests = *tests
	config.Paths = *paths
	config.Packages = flag.Args()

	if !config.Missing && !config.All {
		// fmt.Fprintln(os.Stderr, "Error: one report option is needed.")
		// fmt.Fprintln(os.Stderr, "")

		flag.Usage()
		os.Exit(0)
	}

	if cwd, err := os.Getwd(); err == nil {
		config.WorkingDirectory = cwd
	} else {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if linter, err := vendorlint.NewLinter(config); err == nil {
		linter.Report()
	} else {
		// error!
	}
}
