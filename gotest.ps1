cd utils
go clean -modcache
go mod tidy
go test

cd ../state
go clean -modcache
go mod tidy
go test

cd ../commands
go clean -modcache
go mod tidy
go test

cd ../socketid_client
go clean -modcache
go mod tidy
go test

cd ..
go clean -modcache
go mod tidy

go build .