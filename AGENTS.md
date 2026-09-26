# Agent instructions

## Repository map

- The root module is `go.osspkg.com/do`; its source and tests live at the repository root.
- `mapreduce`, `monad`, `pipline`, `words`, and `workerpool` are separate importable packages.
- `pipline` is the existing package and import path spelling. Preserve it unless a deliberate compatibility change is requested.
- `README.md` is the user-facing overview; keep examples and documented behavior aligned with the exported API.
- For usage questions or changes that use this module, read `skills/go-do-usage/SKILL.md` and the relevant linked reference.

## Go version and commands

Run commands from the repository root. `go.mod` requires Go 1.26.0.

- Run all tests with `go test ./...`.
- For changes to the core helpers, worker pool, or map-reduce package, run `go test . ./workerpool ./mapreduce` first.
- Format changed Go files with `gofmt -w <files>`.
- The GitHub Actions workflow is `.github/workflows/ci.yml`; it invokes `make ci`.
- `make ci` runs `make pre-commit`, which installs the latest `goppy` tool, runs `goppy setup-lib`, and then runs license, lint, test, and build targets. Review those setup and generated-file side effects before running it locally.
- `go.mod` requires Go 1.26.0.

## Change guidance

- Add regression tests for changed behavior. For concurrency changes, cover cancellation, shutdown, and channel ownership; use bounded waits in tests that could hang.
- Keep package boundaries intact. The root package contains generic slice/map helpers, conditionals, numeric helpers, panic recovery, and async helpers; the subpackages provide focused APIs.
- Update `README.md` when installation, public behavior, or usage examples change.
- Avoid changing exported names or import paths without considering downstream compatibility.
- Do not run publication, deployment, or release commands as part of local validation.
