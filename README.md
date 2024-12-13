# proxaudit

proxaudit is a binary that transparently instrument any program for HTTP and HTTPs requests.

## Getting started

```shell
brew tap juliendoutre/proxaudit https://github.com/juliendoutre/proxaudit
brew install proxaudit
mkcert -install
```

## Usage

```shell
proxaudit -- curl http://google.com
proxaudit -- curl https://google.com
proxaudit # Read from stdin
proxaudit -output logs.jsonl -- pip install requests # Write logs to file
proxaudit -server # Run the proxy server only (no command wrapping)
```

## Development

### Lint the code

```shell
brew install golangci-lint hadolint
golangci-lint run
hadolint ./Dockerfile
```

### Release a new version

```shell
brew install goreleaser
goreleaser check syft
git tag -a v0.1.0 -m "First release"
git push origin v0.1.0
rm -rf ./dist
GITHUB_TOKEN=$(gh auth token) goreleaser release
```
