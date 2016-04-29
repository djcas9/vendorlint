package vendorlint

// Config holds application specific
// configuration options.
type Config struct {
	Packages         []string
	Tests            bool
	WorkingDirectory string
}

// NewConfig created a new configuration struct
// and returns with empty defaults.
func NewConfig() *Config {
	return &Config{
		Packages: []string{},
		Tests:    false,
	}
}
