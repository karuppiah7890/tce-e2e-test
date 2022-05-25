check: compilercheck

# check compilation errors for all tool code
# and test code
compilercheck:
	go build -v ./...
	go list ./... | xargs -I {} go test -v -c {}

# TODO: Add target to run the E2E tests

tidy:
	go mod tidy -v -compat=1.17

# TODO: Add target to create standalone binary which can run the E2E tests

# TODO: Add target to create standalone Docker image which can run the TCE E2E tests

# TODO: Add targets to do linting - golangci-lint (staticcheck etc), go.mod and go.sum being up to date,
# custom linting - using internal log package and not std or other log package
