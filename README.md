# proxaudit

## Usage

```shell
brew install mkcert
mkcert -install
go run ./main.go -- curl http://google.com
go run ./main.go -- curl https://google.com
go run ./main.go # Read from stdin
go run ./main.go -output logs.jsonl -- curl https://google.com # Write logs to file
```
