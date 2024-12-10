# proxaudit

## Usage

```shell
brew install mkcert
mkcert -install
go run ./main.go
http_proxy=http://localhost:8000 curl http://google.com
https_proxy=http://localhost:8000 curl https://google.com
```
