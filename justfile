
lint:
  golangci-lint run -E goimports,sqlclosecheck,bodyclose,asciicheck,misspell,errorlint --timeout 5m -e "errors.As" -e "errors.Is" ./...

lint-fix:
	golangci-lint run -E goimports,sqlclosecheck,bodyclose,asciicheck,misspell,errorlint --timeout 5m -e "errors.As" -e "errors.Is" ./... --fix

run cmd: 
  go run cmd/sci/*.go {{cmd}}

tidy: 
  go mod tidy

last-release:
  git describe --tags `git rev-list --tags --max-count=1`

release semver:
    git tag {{semver}}
    git push --tags