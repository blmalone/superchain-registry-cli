
lint:
  golangci-lint run -E goimports,sqlclosecheck,bodyclose,asciicheck,misspell,errorlint --timeout 5m -e "errors.As" -e "errors.Is" ./...

lint-fix:
	golangci-lint run -E goimports,sqlclosecheck,bodyclose,asciicheck,misspell,errorlint --timeout 5m -e "errors.As" -e "errors.Is" ./... --fix

run cmd: 
  go run sure/*.go {{cmd}}

tidy: 
  go mod tidy

last-release:
  @echo "\n#### Remember to update the version in sure/version.go ####\n"
  @git describe --tags `git rev-list --tags --max-count=1`
  @echo "\n"

release semver:
    git tag {{semver}}
    git push --tags