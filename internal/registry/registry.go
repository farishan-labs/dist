package registry

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Channel struct {
	Name    string `toml:"name"`
	Type    string `toml:"type"` // github, gumroad, cloudflare, npm, content
	Enabled bool   `toml:"enabled"`
	// GitHub
	Repo string `toml:"repo"`
	// Gumroad
	ProductID string `toml:"product_id"`
	// Cloudflare
	Project string `toml:"project"`
	// npm
	Package string `toml:"package"`
}

type Registry struct {
	Channels []Channel `toml:"channels"`
}

// Load reads the registry from the given path. If the file does not exist, returns an empty registry.
func Load(path string) (Registry, error) {
	var r Registry
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return r, nil
	}
	if _, err := toml.DecodeFile(path, &r); err != nil {
		return r, err
	}
	return r, nil
}

// DefaultPath returns the default global registry path.
func DefaultPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "fhl", "distribution.toml")
}
