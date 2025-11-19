# Agent Guidelines for this Repository

## Build & Test
- **Build**: `cd go && go build ./...`
- **Test All**: `cd go && go test ./...`
- **Test Single**: `cd go && go test -v -run <TestRegex> <PackagePath>`
  - Example: `cd go && go test -v -run TestFetch ./crawler/reddit`
- **Run App**: `cd go && go run cmd/crawler/main.go`

## Code Style & Conventions
- **Go Version**: Use Go 1.23.4 patterns.
- **Formatting**: Run `gofmt` on all changes. Group imports: stdlib, third-party, internal.
- **Structure**: Follow `cmd/` (entry points) and `internal/` or flat pkg layout as seen.
- **Interfaces**: Define small interfaces (e.g., `Source`, `Notifier`) in `interfaces.go`.
- **Context**: Pass `context.Context` as the first argument to functions performing I/O.
- **Error Handling**: Return errors, don't panic. Use `fmt.Errorf` with wrapping `%w`.
- **Naming**: PascalCase for exported, camelCase for private.
- **Paths**: Always use absolute paths, starting with `/mnt/c/Programming/utilities`.

## Environment
- **Repo Root**: `/mnt/c/Programming/utilities`
- **Working Dir**: Commands usually run from `go/` subdirectory.
