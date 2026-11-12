# dist

Distribution layer for Faris Han Labs: a single command surface for shipping and
measuring across channels (GitHub, Gumroad, Cloudflare, npm, content).

Status: v0 (M1-M4 complete). See `SPEC.md` and `PLAN.md`.

## Build

```sh
go build -o dist ./cmd/dist
```

## Install

```sh
go install github.com/farishan-labs/dist/cmd/dist@latest
```

Or download a release binary from GitHub Releases (built via goreleaser).

## Usage

```sh
dist status                 # channel/product inventory
dist status --json          # machine-readable
dist ship                   # dry-run by default
dist ship --no-dry-run --yes  # execute (approval skipped with --yes)
dist ship --channel workflow-kit
dist metrics                # JSON metrics
dist doctor                 # tools/auth/registry health
dist --version
```

## Registry

Create `~/.config/fhl/distribution.toml`:

```toml
schema_version = 1

[[channels]]
name = "workflow-kit"
type = "github"
enabled = true
repo = "farishan/workflow-kit"

[[channels]]
name = "farishan1-store"
type = "gumroad"
enabled = true
product_id = "all"

[[channels]]
name = "labs-site"
type = "cloudflare"
enabled = true
project = "labs-farishan-com"
```

Per-repo config: create `.fhl-dist.toml` in a repo and load with the registry
package (`LoadRepo`).

## Channels

- `github` — requires `repo`; uses goreleaser if config present else `gh release create`.
- `gumroad` — verifies auth via `gumroad products list`.
- `cloudflare` — requires `project`; uses `cf deploy`.
- `npm` — runs `npm publish`.
- `content` — manual step (dev.to, essays, Product Hunt).

## Metrics

`dist metrics` emits JSON: Gumroad product/sales/revenue counts (via gumroad CLI)
plus optional GoatCounter traffic (set `GOATCOUNTER_URL`, `GOATCOUNTER_TOKEN`).
Views via Gumroad MCP are best-effort (endpoint may require a fresh token).

## Growth-dashboard integration

The growth-dashboard can consume `dist metrics --json` output or call the
registry directly. Keep `growth-dashboard` as the metrics UI; `dist` is the
command surface + event log.

## Event log

Ship actions are appended to `~/.local/share/fhl/dist/events.jsonl` (JSONL).

## Secrets

Use existing Bitwarden/.env. Never commit tokens. `dist doctor` checks auth.

## License

GPL-3.0 (community) per P-006/P-023. See `LICENSE`.
