# PLAN — dist v0

## Milestone 0 — Scaffold (done)
## Milestone 1 — Registry + status (done)
- Channel-registry TOML schema v1 (global + per-repo via LoadRepo/Save).
- `dist status` with --json, validation, enabled filter.
- Tests: registry parse, validation, save/load, enabled, find, loadrepo.
## Milestone 2 — Ship orchestration (done)
- `dist ship` with --dry-run (default), --yes, --channel.
- Approval gate on real runs; subprocess delegation (gh/goreleaser/cf/gumroad/npm).
- JSONL event log at ~/.local/share/fhl/dist/events.jsonl.
## Milestone 3 — Metrics (done)
- `dist metrics --json`: Gumroad products/sales/revenue via gumroad CLI; optional GoatCounter; best-effort MCP views.
- Growth-dashboard integration documented in README.
## Milestone 4 — Polish (done)
- `dist doctor`: tools/auth/registry health.
- goreleaser config (.goreleaser.yaml), README, install docs.
- Tag v0.1.0 pending.
- Go module, cobra CLI skeleton, SPEC/PLAN, README, LICENSE (GPL-3.0).

## Milestone 1 — Registry + status (week 1)
- Define channel-registry TOML schema (global + per-repo).
- Implement `dist status` reading registry, printing channel/product inventory.
- Tests for registry parse + status output.

## Milestone 2 — Ship orchestration (week 2)
- Implement `dist ship` with dry-run default and approval gate.
- Delegate to gh/goreleaser/wrangler/gumroad via subprocess.
- Event log (JSONL) of ship actions.

## Milestone 3 — Metrics (week 3)
- Implement `dist metrics --json` (Gumroad analytics + GoatCounter API).
- Document growth-dashboard integration (consume JSON).

## Milestone 4 — Polish (week 4)
- Doctor checks (auth, tool versions).
- README, install, goreleaser config.
- Tag v0.1.0.

## Verification
- `go build ./...`, `go test ./...`, `go vet ./...`.
- Manual dry-run ship on one repo.

## Milestone 5 — Content channel (planned)
- dev.to article: "dist: one command surface for shipping and measuring across channels" (install, status, ship dry-run, metrics).
- Product Hunt launch: tagline, gallery, demo GIF, install one-liner.
- essays.farishan.com deep-dive: registry schema + event log + metrics loop.
- npm/design-tokens cross-link and Cloudflare Pages demo site.
- Checklist: README quickstart, `dist --version`, `dist doctor` green, sample `dist metrics --json`, release notes for v0.3.0.
