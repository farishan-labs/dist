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
	content := "schema_version = 1\n\n[[channels]]\nname = \"wk\"\ntype = \"github\"\nenabled = true\nrepo = \"farishan/workflow-kit\"\n"
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
	if r.SchemaVersion != 1 {
		t.Fatalf("schema version: %d", r.SchemaVersion)
	}
	if err := r.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestValidateDuplicateName(t *testing.T) {
	r := Registry{SchemaVersion: 1, Channels: []Channel{
		{Name: "a", Type: "github", Repo: "x/y", Enabled: true},
		{Name: "a", Type: "github", Repo: "x/z", Enabled: true},
	}}
	if err := r.Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestValidateUnknownType(t *testing.T) {
	r := Registry{SchemaVersion: 1, Channels: []Channel{
		{Name: "a", Type: "bogus", Enabled: true},
	}}
	if err := r.Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestValidateMissingRepo(t *testing.T) {
	r := Registry{SchemaVersion: 1, Channels: []Channel{
		{Name: "a", Type: "github", Enabled: true},
	}}
	if err := r.Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestValidateMissingCloudflareProject(t *testing.T) {
	r := Registry{SchemaVersion: 1, Channels: []Channel{
		{Name: "a", Type: "cloudflare", Enabled: true},
	}}
	if err := r.Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestValidateMissingNpmPackage(t *testing.T) {
	r := Registry{SchemaVersion: 1, Channels: []Channel{
		{Name: "a", Type: "npm", Enabled: true},
	}}
	if err := r.Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestValidateMissingContentPlatform(t *testing.T) {
	r := Registry{SchemaVersion: 1, Channels: []Channel{
		{Name: "a", Type: "content", Enabled: true},
	}}
	if err := r.Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "sub", "distribution.toml")
	r := Registry{SchemaVersion: 1, Channels: []Channel{
		{Name: "wk", Type: "github", Enabled: true, Repo: "farishan/workflow-kit"},
	}}
	if err := Save(p, r); err != nil {
		t.Fatal(err)
	}
	r2, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(r2.Channels) != 1 || r2.Channels[0].Name != "wk" {
		t.Fatalf("roundtrip failed: %+v", r2)
	}
}

func TestSaveInvalidFails(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "distribution.toml")
	r := Registry{SchemaVersion: 1, Channels: []Channel{
		{Name: "", Type: "github", Repo: "x/y"},
	}}
	if err := Save(p, r); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestEnabled(t *testing.T) {
	r := Registry{Channels: []Channel{
		{Name: "a", Type: "github", Repo: "x/y", Enabled: true},
		{Name: "b", Type: "github", Repo: "x/z", Enabled: false},
	}}
	if got := len(r.Enabled()); got != 1 {
		t.Fatalf("expected 1 enabled, got %d", got)
	}
}

func TestFind(t *testing.T) {
	r := Registry{Channels: []Channel{
		{Name: "a", Type: "github", Repo: "x/y"},
	}}
	c, idx := r.Find("a")
	if idx != 0 || c.Name != "a" {
		t.Fatalf("find failed: %+v %d", c, idx)
	}
	_, idx2 := r.Find("nope")
	if idx2 != -1 {
		t.Fatalf("expected -1, got %d", idx2)
	}
}

func TestLoadRepo(t *testing.T) {
	dir := t.TempDir()
	content := "schema_version = 1\n\n[[channels]]\nname = \"repo-gh\"\ntype = \"github\"\nenabled = true\nrepo = \"farishan/repo\"\n"
	if err := os.WriteFile(filepath.Join(dir, ".fhl-dist.toml"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	r, err := LoadRepo(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Channels) != 1 || r.Channels[0].Name != "repo-gh" {
		t.Fatalf("loadrepo failed: %+v", r)
	}
}
