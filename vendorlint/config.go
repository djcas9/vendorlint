package vendorlint

// Config holds application specific
// configuration options.
type Config struct {
	All              bool
	Missing          bool
	Packages         []string
	Tests            bool
	WorkingDirectory string
}

// NewConfig created a new configuration struct
// and returns with empty defaults.
func NewConfig() *Config {
	return &Config{
		All:      false,
		Missing:  false,
		Packages: []string{},
		Tests:    false,
	}
}
