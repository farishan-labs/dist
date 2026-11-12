# SPEC — dist (FHL distribution layer)

## Purpose
Single command surface for shipping and measuring across all FHL distribution channels.
Orchestrates existing wheels; does not reimplement them (P-001, P-007).

## Principles
P-001 reuse, P-005 single static binary (Go, no CGO), P-006 open/shared (GPL-3.0 default),
P-009 one tool one job (money path), P-010 human sets direction (dry-run default, approval gate).

## v0 scope
1. `dist status` — inventory of channels/products and last-shipped state.
2. `dist ship` — orchestrate GitHub release (goreleaser) + Cloudflare Pages (wrangler) +
   Gumroad product sync (MCP/API). Human approval gate + dry-run default.
3. `dist metrics --json` — pull Gumroad sales/revenue + GoatCounter views; print JSON.
   Metrics UI lives in growth-dashboard (extended), not here.
4. Canonical JSON event log (P-012) as interchange; no new DB.
5. Secrets via existing Bitwarden/.env (P-015/P-016); no account creation (P-019).

## Deferred (v1+)
npm/Homebrew publish, content (dev.to/PH), audience capture, Tauri UI.
ESP selection supported via `esp` channel type; UTM attribution join included in `dist metrics`.

## Channel registry
Single source of truth: `~/.config/fhl/distribution.toml` (global) and/or per-repo
`.fhl-dist.toml`. Schema v1 defined in `internal/registry` (global + per-repo).

## Non-goals
No marketing automation, no auto-publishing, no reinvention of gh/goreleaser/cf/npm/Gumroad.

## License
GPL-3.0 (community) per P-006/P-023.
