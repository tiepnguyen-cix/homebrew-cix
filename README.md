# cix

Run and debug GitLab CI pipelines locally. No push required.

```
cix run --job build
```

---

## The problem

Every GitLab CI change requires a push. You wait 5–10 minutes for the pipeline, it fails, you fix, push again. Each cycle burns CI minutes and kills your flow.

`cix` runs your `.gitlab-ci.yml` jobs directly on your machine inside Docker — same environment, same scripts, zero waiting.

## Install

**macOS (Homebrew)**

```bash
brew tap tiepnguyen-cix/cix
brew install cix
```

**Linux / macOS (curl)**

```bash
curl -sSL https://github.com/tiepnguyen-cix/cix/releases/latest/download/install.sh | sh
```

**Go**

```bash
go install github.com/tiepnguyen-cix/cix/cmd/cix@latest
```

## Usage

```bash
# List all jobs in the pipeline
cix list

# Validate .gitlab-ci.yml syntax
cix validate

# Run a specific job
cix run --job build

# Run all jobs in a stage
cix run --stage test

# Preview what would run without executing
cix run --job build --dry-run

# Keep container alive after failure for debugging
cix run --job build --keep-on-fail

# Show full command output
cix run --job build --verbose

# Use a different CI file
cix run --job build -f path/to/.gitlab-ci.yml
```

## Demo

```
$ cix run --job build

cix v0.1.0 · GitLab CI local runner

▶ job: build  (node:20-alpine)
  ✓ npm ci                                        12.3s
  ✓ npm run build                                  8.1s
  ✗ npm run test

FAILED · build exited with code 1

stderr output:
  Error: Cannot find module './config'
  at Object.<anonymous> (/app/src/index.js:3:18)

tip: run with --keep-on-fail to inspect the container
```

## Local config

Create a `.cix.yml` file next to your `.gitlab-ci.yml` to override variables for local runs:

```yaml
# .cix.yml
variables:
  NODE_ENV: development
  DATABASE_URL: postgres://localhost/mydb

secrets:
  API_KEY: $API_KEY  # resolved from shell environment

docker:
  pull_policy: if-not-present
  network: bridge
```

Add `.cix.yml` to your `.gitignore` — it is for local use only.

## Supported GitLab CI keywords

| Keyword | Status |
|---|---|
| `image` | ✅ |
| `stage` | ✅ |
| `script` | ✅ |
| `before_script` | ✅ |
| `after_script` | ✅ |
| `variables` | ✅ |
| `artifacts.paths` | ✅ |
| `needs` | ✅ |
| `allow_failure` | ✅ |
| `extends` | 🔜 v0.2 |
| `include` | 🔜 v0.2 |
| `services` | 🔜 v0.2 |
| `rules` | 🔜 v0.2 |
| `cache` | 🔜 v0.2 |

## Predefined CI variables

`cix` automatically injects these variables so your scripts work as expected:

| Variable | Value |
|---|---|
| `CI` | `true` |
| `CI_JOB_NAME` | current job name |
| `CI_COMMIT_BRANCH` | current git branch |
| `CI_COMMIT_SHORT_SHA` | current git commit SHA |
| `CI_PIPELINE_SOURCE` | `local` |
| `CI_PROJECT_DIR` | `/builds/project` |

## Requirements

- [Docker](https://docs.docker.com/get-docker/) running locally
- Go 1.22+ (only if installing via `go install`)

## vs act

[act](https://github.com/nektos/act) is a great tool for running GitHub Actions locally. `cix` focuses on GitLab CI:

| | cix | act |
|---|---|---|
| GitLab CI | ✅ | ❌ |
| GitHub Actions | 🔜 | ✅ |
| Step-by-step output | ✅ | ❌ |
| `--keep-on-fail` | ✅ | ❌ |
| No pipeline rewrite | ✅ | ✅ |
| Single binary | ✅ | ✅ |

## Contributing

```bash
git clone https://github.com/tiepnguyen-cix/cix
cd cix
go mod tidy
go test ./...
```

Pull requests are welcome.

## License

MIT