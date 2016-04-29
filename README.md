# vendorlint

vendorlint was designed to help me prevent imports that have not been added to the vendor directory. 
vendorlint will also output the source file line/column for the missing dependency.

## Usage

  ```bash
Usage: vendorlint [-t] package [package ...]
  -t	include test dependencies
  -v	vendorlint version number
  ```

## Example Output

  ```
dweb@dweb-linux vendorlint (master*) $ ./bin/vendorlint ./...
[X] Dependency not vendored: github.com/fatih/color
  * vendorlint/lint.go:14:2
[X] Dependency not vendored: github.com/kisielk/gotool
  * vendorlint/lint.go:15:2
  ```

## TODO

  * Add extra rules around vendor best practices.
  * Test logic need to be completed
  * Add tests

## Acknowledgments

  I was inspired by [vendorcheck](https://github.com/FiloSottile/vendorcheck) but wanted to add extra functionality.
  It became clear that what I wanted to add would result a complete rewrite. I decide to just start something new.

## Self-Promotion

  * [blog](http://dweb.io/)
  * [Twitter](http://twitter.com/mephux)
