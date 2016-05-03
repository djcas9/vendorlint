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

	// BasePackageWords holds the cwd base paths
	BasePackageWords = []string{}

	// PackagePaths holds the abs package paths
	PackagePaths = []string{}
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
		if !hasVendorMatch(i.Name) {
			color.Red("[X] Dependency not vendored: %s\n", i.Name)
			fmt.Fprintf(os.Stderr, "  * %s\n", i.Position.String())
		}
	}
}

func hasVendorMatch(name string) bool {
	for _, fp := range PackagePaths {
		if _, err := os.Stat(path.Join(fp,
			"vendor", name)); err != nil {
			return false
		}
	}

	return true
}

func (l *Linter) run() error {

	var files []string
	argPaths := gotool.ImportPaths(flag.Args())

	for _, imp := range argPaths {

		baseName := filepath.Base(imp)
		fullPath, _ := filepath.Abs(imp)

		if baseName == "." {
			BasePackageWords = append(BasePackageWords, filepath.Base(l.Config.WorkingDirectory))
			PackagePaths = append(PackagePaths, l.Config.WorkingDirectory)
		} else {
			BasePackageWords = append(BasePackageWords, baseName)
			PackagePaths = append(PackagePaths, fullPath)
		}

		filepath.Walk(imp, func(path string,
			f os.FileInfo, err error) error {

			if strings.Contains(path, ".go") {
				if strings.Contains(path, "_test.go") {
					if l.Config.Tests {
						files = append(files, path)
					}
				} else {
					files = append(files, path)
				}
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

				if !hasIgnoreKeyboard(name, IgnoreKeywords) &&
					!hasIgnoreKeyboard(position.String(), IgnoreKeywords) {

					if !hasIgnoreKeyboard(name, BasePackageWords) {

						imports = append(imports, Import{
							Filename: filepath.Base(pt),
							Name:     name,
							Position: position,
						})

					}

				}
			}

		}
	}

	l.Imports = imports

	return nil
}

func hasIgnoreKeyboard(p string, list []string) bool {
	for _, i := range list {
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
