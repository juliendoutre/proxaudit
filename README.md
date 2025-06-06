# proxaudit

proxaudit is a binary that transparently instrument any program for HTTP and HTTPs requests.

## Getting started

```shell
brew tap juliendoutre/proxaudit https://github.com/juliendoutre/proxaudit
brew install mkcert proxaudit
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

Check out [examples/k8s/](./examples/k8s/) for an example of using proxaudit as a sidecar container in Kubernetes.

## Development

### Lint the code

```shell
brew install golangci-lint hadolint
golangci-lint run
hadolint ./Dockerfile
```

### Release a new version

```shell
git tag -a v0.1.0 -m "New release"
git push origin v0.1.0
```
