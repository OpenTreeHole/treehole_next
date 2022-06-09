rm -rf ./tests/data
go clean -testcache
cls
go test ./tests -v