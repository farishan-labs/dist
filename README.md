# dist

Distribution layer for Faris Han Labs: a single command surface for shipping and
measuring across channels (GitHub, Gumroad, Cloudflare, npm, content).

Status: v0 scaffold. See `SPEC.md` and `PLAN.md`.

## Build

```sh
go build -o dist ./cmd/dist
```

## Usage

```sh
dist status
dist ship --dry-run
dist metrics
```

## Registry

Create `~/.config/fhl/distribution.toml`:

```toml
[[channels]]
name = "workflow-kit"
type = "github"
enabled = true
repo = "farishan/workflow-kit"
```

## License

GPL-3.0 (community) per P-006/P-023. See `LICENSE`.
