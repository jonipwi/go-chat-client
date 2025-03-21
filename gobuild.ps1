go get github.com/googollee/go-socket.io
go mod tidy
go mod download
go clean -modcache
go build .