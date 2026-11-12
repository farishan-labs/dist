# PLAN — dist v0

## Milestone 0 — Scaffold (done)
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
