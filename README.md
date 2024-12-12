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
proxaudit -output logs.jsonl -- curl https://google.com # Write logs to file
```

## Development

### Lint the code

```shell
brew install golangci-lint
golangci-lint run
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
