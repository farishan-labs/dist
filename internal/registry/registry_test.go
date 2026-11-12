package registry

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMissingFile(t *testing.T) {
	r, err := Load(filepath.Join(t.TempDir(), "nope.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Channels) != 0 {
		t.Fatalf("expected empty, got %d", len(r.Channels))
	}
}

func TestLoadValidFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "distribution.toml")
	content := `[[channels]]
name = "wk"
type = "github"
enabled = true
repo = "farishan/workflow-kit"
`
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	r, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Channels) != 1 {
		t.Fatalf("expected 1, got %d", len(r.Channels))
	}
	if r.Channels[0].Name != "wk" || !r.Channels[0].Enabled {
		t.Fatalf("bad channel: %+v", r.Channels[0])
	}
}
