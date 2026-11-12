package registry

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const SchemaVersion = 1

const (
	TypeGitHub     = "github"
	TypeGumroad    = "gumroad"
	TypeCloudflare = "cloudflare"
	TypeNPM        = "npm"
	TypeContent    = "content"
)

type Channel struct {
	Name        string `toml:"name"`
	Type        string `toml:"type"`
	Enabled     bool   `toml:"enabled"`
	Repo        string `toml:"repo,omitempty"`
	ReleaseTool string `toml:"release_tool,omitempty"`
	ProductID   string `toml:"product_id,omitempty"`
	Project     string `toml:"project,omitempty"`
	Package     string `toml:"package,omitempty"`
	Platform    string `toml:"platform,omitempty"`
	LastShipped string `toml:"last_shipped,omitempty"`
	Notes       string `toml:"notes,omitempty"`

	// execution
	Dir string `toml:"dir,omitempty"`
}

type Registry struct {
	SchemaVersion int       `toml:"schema_version"`
	Channels      []Channel `toml:"channels"`
}

var validTypes = map[string]bool{
	TypeGitHub: true, TypeGumroad: true, TypeCloudflare: true, TypeNPM: true, TypeContent: true,
}

func (r Registry) Validate() error {
	if r.SchemaVersion != 0 && r.SchemaVersion != SchemaVersion {
		return fmt.Errorf("unsupported schema_version %d (want %d)", r.SchemaVersion, SchemaVersion)
	}
	seen := map[string]bool{}
	for i, c := range r.Channels {
		if c.Name == "" {
			return fmt.Errorf("channel %d: name required", i)
		}
		if seen[c.Name] {
			return fmt.Errorf("duplicate channel name %q", c.Name)
		}
		seen[c.Name] = true
		if !validTypes[c.Type] {
			return fmt.Errorf("channel %q: unknown type %q", c.Name, c.Type)
		}
		switch c.Type {
		case TypeGitHub:
			if c.Repo == "" {
				return fmt.Errorf("channel %q: repo required for github", c.Name)
			}
		case TypeCloudflare:
			if c.Project == "" {
				return fmt.Errorf("channel %q: project required for cloudflare", c.Name)
			}
		case TypeNPM:
			if c.Package == "" {
				return fmt.Errorf("channel %q: package required for npm", c.Name)
			}
		case TypeContent:
			if c.Platform == "" {
				return fmt.Errorf("channel %q: platform required for content", c.Name)
			}
		}
	}
	return nil
}

func Load(path string) (Registry, error) {
	var r Registry
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return r, nil
	}
	if _, err := toml.DecodeFile(path, &r); err != nil {
		return r, err
	}
	if r.SchemaVersion == 0 {
		r.SchemaVersion = SchemaVersion
	}
	return r, nil
}

func LoadRepo(dir string) (Registry, error) {
	return Load(filepath.Join(dir, ".fhl-dist.toml"))
}

func Save(path string, r Registry) error {
	if err := r.Validate(); err != nil {
		return err
	}
	if r.SchemaVersion == 0 {
		r.SchemaVersion = SchemaVersion
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return toml.NewEncoder(f).Encode(r)
}

func DefaultPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "fhl", "distribution.toml")
}

func (r Registry) Enabled() []Channel {
	out := []Channel{}
	for _, c := range r.Channels {
		if c.Enabled {
			out = append(out, c)
		}
	}
	return out
}

func (r Registry) Find(name string) (Channel, int) {
	for i, c := range r.Channels {
		if c.Name == name {
			return c, i
		}
	}
	return Channel{}, -1
}

var ErrNotFound = errors.New("channel not found")
