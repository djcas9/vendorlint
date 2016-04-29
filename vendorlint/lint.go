package vendorlint

import (
	"flag"
	"fmt"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/kisielk/gotool"
)

var (
	// IgnoreKeywords are used to prevent
	// imports that are not required for the main
	// package
	IgnoreKeywords = []string{
		"vendor/",
		"Godeps/",
	}
)

// Import holds metadata useful for reporting
type Import struct {
	Name     string
	Filename string
	Position token.Position
}

// Linter holds the application configuration
// file, import reuslts and the context struct used
// by go/build
type Linter struct {
	Config  *Config
	Imports []Import
	Context build.Context
}

// NewLinter returns a Linter struct using the
// passed application configurations
func NewLinter(config *Config) (*Linter, error) {
	linter := &Linter{
		Config: config,
	}

	err := linter.run()

	return linter, err
}

// Report writes all failed vendor deps
func (l *Linter) Report() {
	for _, i := range l.Imports {

		if strings.Contains(i.Name,
			filepath.Base(l.Config.WorkingDirectory)) {
			continue
		}

		if _, err := os.Stat(path.Join(l.Config.WorkingDirectory,
			"vendor", i.Name)); err != nil {
			color.Red("[X] Dependency not vendored: %s\n", i.Name)
			fmt.Fprintf(os.Stderr, "  * %s\n", i.Position.String())
		}
	}
}

func (l *Linter) run() error {

	var files []string
	for _, imp := range gotool.ImportPaths(flag.Args()) {
		filepath.Walk(imp, func(path string,
			f os.FileInfo, err error) error {

			if strings.Contains(path, ".go") {
				files = append(files, path)
			}

			return err
		})
	}

	var imports []Import
	cache := make(map[string]bool)

	for _, pt := range files {
		fset := token.NewFileSet()
		astFile, _ := parser.ParseFile(fset, pt, nil, parser.ImportsOnly)

		for _, ii := range astFile.Imports {
			position := fset.File(astFile.Pos()).Position(ii.Pos())
			name := strings.Replace(ii.Path.Value, "\"", "", -1)

			if cache[name] {
				continue
			}

			cache[name] = true

			if !strings.HasPrefix(name, build.Default.GOROOT) &&
				!isStandardImportPath(ii.Path.Value) {

				if !hasIgnoreKeyboard(name) &&
					!hasIgnoreKeyboard(position.String()) {

					imports = append(imports, Import{
						Filename: filepath.Base(pt),
						Name:     name,
						Position: position,
					})
				}
			}

		}
	}

	l.Imports = imports

	return nil
}

func hasIgnoreKeyboard(p string) bool {
	for _, i := range IgnoreKeywords {
		if strings.Contains(p, i) {
			return true
		}
	}

	return false
}

// from https://github.com/golang/go/blob/87bca88/src/cmd/go/pkg.go#L183-L194
func isStandardImportPath(path string) bool {
	i := strings.Index(path, "/")
	if i < 0 {
		i = len(path)
	}
	elem := path[:i]
	return !strings.Contains(elem, ".")
}
