# omm

Keyboard-driven TUI task manager for the terminal, built in Go with Bubbletea.

## Common Commands

Use **just** (preferred over raw Go commands):

| Command        | Alias     | Action              |
|----------------|-----------|---------------------|
| `just run`     | `just r`  | `go run .`          |
| `just build`   | `just b`  | `go build .`        |
| `just test`    | `just t`  | `go test ./...`     |
| `just fmt`     | `just f`  | `gofumpt -l -w .`   |
| `just lint`    | `just l`  | `golangci-lint run` |
| `just install` | `just i`  | `go install .`      |
| `just tidy`    | `just ti` | `go mod tidy`       |
| `just vuln`    | `just v`  | `govulncheck ./...` |

## Architecture

- `cmd/` -- CLI layer (Cobra + Viper). Config precedence: flags > env (`OMM_*`) > TOML config > defaults
- `internal/ui/` -- Bubbletea TUI (Elm Architecture: Model/Update/View). Messages in `msgs.go`, async commands in `cmds.go`
- `internal/ui/theme/` -- Built-in theme registry (8 themes). Each theme defines a semantic color palette used across the TUI and markdown rendering
- `internal/persistence/` -- SQLite data layer (pure-Go driver). Custom append-only migration system in `migrations.go`
- `internal/types/` -- Domain types (`Task`, `TaskPrefix`, `ContextBookmark`)
- `internal/utils/` -- Shared utilities

## Key Conventions

- Format with `gofumpt`, lint with `golangci-lint` v2 (config in `.golangci.yml`)
- Tests use `testing` + `testify` (`assert`/`require`). Table-driven test style preferred
- Persistence tests use in-memory SQLite via `TestMain` setup
- Error variables: `errCamelCase` (e.g., `errCouldntGetHomeDir`)
- Bubbletea messages: `*Msg` suffix. Commands return `tea.Cmd`
- Assets embedded via `//go:embed` (guide content, help text, changelog)
- Max 10,000 active tasks; context capped at 1 MB
- Built-in theme system (`internal/ui/theme/`); styles in `styles.go` are constructed from the active theme

## Release

- GoReleaser v2 (`.goreleaser.yaml`): linux + darwin, amd64 + arm64
- Artifacts signed with cosign (Sigstore)
- Releases created as drafts, published to Homebrew tap `dhth/tap/omm`
- Changelog follows Keep a Changelog format with SemVer
