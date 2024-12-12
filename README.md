# proxaudit

proxaudit is a binary that transparently instrument any program for HTTP and HTTPs requests.

## Usage

```shell
brew install mkcert
mkcert -install
go run ./main.go -- curl http://google.com
go run ./main.go -- curl https://google.com
go run ./main.go # Read from stdin
go run ./main.go -output logs.jsonl -- curl https://google.com # Write logs to file
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
