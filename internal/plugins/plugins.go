package plugins

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Plugin struct {
	Name    string   `toml:"name"`
	Binary  string   `toml:"binary"`
	Args    []string `toml:"args,omitempty"`
	Enabled bool     `toml:"enabled"`
}

type Manifest struct {
	Plugins []Plugin `toml:"plugin"`
}

type Envelope struct {
	Metrics json.RawMessage `json:"metrics"`
	Notes   []string        `json:"notes"`
}

func Path() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "fhl", "dist-plugins.toml")
}

func defaults() Manifest {
	return Manifest{Plugins: []Plugin{
		{Name: "gumroad", Binary: "dist-plugin", Args: []string{"gumroad"}, Enabled: true},
		{Name: "goatcounter", Binary: "dist-plugin", Args: []string{"goatcounter"}, Enabled: true},
	}}
}

func Load() (Manifest, error) {
	var m Manifest
	p := Path()
	if _, err := os.Stat(p); os.IsNotExist(err) {
		m = defaults()
		if err := Save(m); err != nil {
			return m, err
		}
		return m, nil
	}
	if _, err := toml.DecodeFile(p, &m); err != nil {
		return m, err
	}
	if len(m.Plugins) == 0 {
		m = defaults()
		_ = Save(m)
	}
	return m, nil
}

func Save(m Manifest) error {
	if err := os.MkdirAll(filepath.Dir(Path()), 0o755); err != nil {
		return err
	}
	f, err := os.Create(Path())
	if err != nil {
		return err
	}
	defer f.Close()
	return toml.NewEncoder(f).Encode(m)
}

func Find(m Manifest, name string) (Plugin, int) {
	for i, p := range m.Plugins {
		if p.Name == name {
			return p, i
		}
	}
	return Plugin{}, -1
}

func Add(name, binary string, args []string) error {
	m, err := Load()
	if err != nil {
		return err
	}
	if _, idx := Find(m, name); idx >= 0 {
		return fmt.Errorf("plugin %s already exists", name)
	}
	m.Plugins = append(m.Plugins, Plugin{Name: name, Binary: binary, Args: args, Enabled: true})
	return Save(m)
}

func Remove(name string) error {
	m, err := Load()
	if err != nil {
		return err
	}
	_, idx := Find(m, name)
	if idx < 0 {
		return fmt.Errorf("plugin %s not found", name)
	}
	m.Plugins = append(m.Plugins[:idx], m.Plugins[idx+1:]...)
	return Save(m)
}

func SetEnabled(name string, enabled bool) error {
	m, err := Load()
	if err != nil {
		return err
	}
	_, idx := Find(m, name)
	if idx < 0 {
		return fmt.Errorf("plugin %s not found", name)
	}
	m.Plugins[idx].Enabled = enabled
	return Save(m)
}

func Run(p Plugin) ([]byte, error) {
	bin := p.Binary
	if bin == "" {
		bin = p.Name
	}
	if _, err := exec.LookPath(bin); err != nil {
		return nil, fmt.Errorf("binary not found: %s", bin)
	}
	args := p.Args
	cmd := exec.Command(bin, args...)
	out, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return out, fmt.Errorf("%w: %s", err, string(ee.Stderr))
		}
		return out, err
	}
	// validate envelope
	var env Envelope
	if err := json.Unmarshal(out, &env); err != nil {
		return out, fmt.Errorf("plugin %s returned invalid envelope: %w", p.Name, err)
	}
	return out, nil
}
