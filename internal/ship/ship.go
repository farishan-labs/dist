package ship

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/farishan-labs/dist/internal/registry"
)

type Result struct {
	Channel string    `json:"channel"`
	Type    string    `json:"type"`
	DryRun  bool      `json:"dry_run"`
	Skipped bool      `json:"skipped"`
	Summary string    `json:"summary"`
	Command string    `json:"command"`
	Output  string    `json:"output"`
	Err     error     `json:"-"`
	Error   string    `json:"error,omitempty"`
	TS      time.Time `json:"ts"`
}

type Event = Result

type Options struct {
	Tag  string
	Dir  string
	File string
}

func toolExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func runCmd(name string, args []string, dir string, dryRun bool) (string, error) {
	if dryRun {
		return fmt.Sprintf("[dry-run] %s %s", name, joinArgs(args)), nil
	}
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func joinArgs(a []string) string {
	out := ""
	for i, s := range a {
		if i > 0 {
			out += " "
		}
		out += s
	}
	return out
}

func planFor(c registry.Channel) (string, []string, string) {
	// returns (binary, args, summary)
	switch c.Type {
	case registry.TypeGitHub:
		repoDir := c.Repo
		if c.Dir != "" {
			repoDir = c.Dir
		}
		if toolExists("goreleaser") {
			cfg := ""
			if repoDir != "" {
				for _, n := range []string{".goreleaser.yaml", ".goreleaser.yml", "goreleaser.yaml"} {
					if _, err := os.Stat(filepath.Join(repoDir, n)); err == nil {
						cfg = n
						break
					}
				}
			}
			if cfg != "" {
				return "goreleaser", []string{"release", "--clean"}, "goreleaser release"
			}
		}
		tag := c.Tag
		if tag == "" && repoDir != "" {
			// try git describe
			out, err := exec.Command("git", "-C", repoDir, "describe", "--tags", "--abbrev=0").Output()
			if err == nil {
				tag = strings.TrimSpace(string(out))
			}
		}
		if tag == "" {
			return "", nil, "no tag; pass --tag or tag the repo"
		}
		return "gh", []string{"release", "create", tag, "--generate-notes"}, "gh release create " + tag
	case registry.TypeCloudflare:
		if toolExists("cf") {
			args := []string{"deploy"}
			if c.Project != "" {
				args = append(args, "--project-name", c.Project)
			}
			if c.Dir != "" {
				args = append(args, c.Dir)
			}
			return "cf", args, "cf deploy"
		}
		if toolExists("wrangler") {
			return "wrangler", []string{"pages", "deploy"}, "wrangler pages deploy"
		}
		return "", nil, "no cloudflare tool found"
	case registry.TypeGumroad:
		if c.File != "" && c.ProductID != "" && c.ProductID != "all" {
			return "gumroad", []string{"products", "update", c.ProductID, "--file", c.File}, "gumroad products update (attach file)"
		}
		return "gumroad", []string{"products", "list", "--json"}, "gumroad products list (verify auth)"
	case registry.TypeNPM:
		return "npm", []string{"publish"}, "npm publish"
	case registry.TypeContent:
		return "", nil, "manual: publish content to " + c.Platform
	case registry.TypeESP:
		return "", nil, "manual: configure ESP " + c.ESP
	}
	return "", nil, "unknown channel type"
}

func Run(channels []registry.Channel, dryRun bool, opts Options) []Result {
	results := []Result{}
	for _, c := range channels {
		if opts.Tag != "" && c.Tag == "" {
			c.Tag = opts.Tag
		}
		if opts.Dir != "" && c.Dir == "" {
			c.Dir = opts.Dir
		}
		if opts.File != "" && c.File == "" {
			c.File = opts.File
		}
		r := Result{Channel: c.Name, Type: c.Type, DryRun: dryRun, TS: time.Now().UTC()}
		if c.Type == registry.TypeContent {
			r.Skipped = true
			r.Summary = "manual step: publish to " + c.Platform
			results = append(results, r)
			continue
		}
		if c.Type == registry.TypeESP {
			r.Skipped = true
			r.Summary = "manual step: configure ESP " + c.ESP
			results = append(results, r)
			continue
		}
		bin, args, summary := planFor(c)
		r.Summary = summary
		r.Command = bin + " " + joinArgs(args)
		if bin == "" {
			r.Skipped = true
			r.Err = fmt.Errorf("no tool available")
			results = append(results, r)
			continue
		}
		if !toolExists(bin) {
			r.Skipped = true
			r.Err = fmt.Errorf("binary not found: %s", bin)
			results = append(results, r)
			continue
		}
		dir := c.Dir
		out, err := runCmd(bin, args, dir, dryRun)
		r.Output = out
		r.Err = err
		if err != nil {
			r.Error = err.Error()
		}
		results = append(results, r)
	}
	return results
}

func EventLogPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", "fhl", "dist", "events.jsonl")
}

func AppendEvents(results []Result) error {
	p := EventLogPath()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	for _, r := range results {
		if err := enc.Encode(r); err != nil {
			return err
		}
	}
	return nil
}
