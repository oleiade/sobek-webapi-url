# AGENTS.md

## Project knowledge

- **Tech stack**: Go 1.25, Sobek (https://github.com/grafana/sobek), Web Platform Test Suite (https://github.com/web-platform-tests/wpt).
- **Goal**: this project implements support for the URL and URLSearchParams WebAPI specification for sobek, with the secondary intent to make it available to k6's sobek-based runtime.
- **File structure**:
  - `url/` - Sobek go module source code (you START from here).
  - `wpt/` - curated Web Platform Test Suite tests that this implementation needs to pass (defined by `wpt.json`).
  - `patches/` - git patches applied to the Web Platform Test Suite in order to be compatible with the sobek javascript runtime.

## Dev environment tips

- Use `go mod tidy` to synchronize the go package's dependencies.
- Use `wptsync` to ensure the web-platform-tests suite is up to date with the `wpt.json` configuration.
- Use `go build ./...` to build the project.

## Testing instructions

- Find the CI plan in the `.github/workflows` folder.
- Run `go test ./...` from the project's root to run every test defined for the package.
- Add or update tests for the code you change, even if nobody asked.

## Code quality

- Use `go vet` to examine the source code and report suspicious constructs
- After writing or modifying code, moving files around or changing imports, run `golangci-lint run` from the project's root to be sure golangci-lint rules pass.

## Documentation practices

- Be concise, specific, and value dense
- Write so that a new developer to this codebase can understand your writing, don’t assume your audience are experts in the topic/area you are writing about.