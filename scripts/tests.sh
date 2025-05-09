#go test ../tests... -v
go test -coverprofile=coverage.out ../pkg/i18n/...
go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html
