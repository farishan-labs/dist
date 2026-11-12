# Release-need matrix (FHL repos)

Derived from git tags and CI workflows (2026-11-12).

| repo | needs releases | mechanism | notes |
|---|---|---|---|
| workflow-kit | yes | GitHub Actions release.yml on v* tag | zip bundle; goreleaser not used |
| monokit | yes (npm) | npm publish | v0.1.0, v0.2 tagged |
| design-tokens | yes (npm) | npm publish | DTCG tokens package |
| agent-tools | maybe (npm) | npm (sweep core extracted) | no tags yet |
| cockpit | not yet | goreleaser config present | M4 polish pending |
| dist | yes | goreleaser (.goreleaser.yaml) | this repo |
| growth-dashboard | no | local app | not distributed |
| essays | no | content | manual publishing |

Commit-message discipline: conventional commits (feat/fix/chore/docs/ci/refactor/test)
observed across repos → compatible with semrel if auto-versioning desired.
