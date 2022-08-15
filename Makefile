NAMESPACE = `echo monosvc`
BUILD_TIME = `date +%FT%T%z`
BUILD_VERSION = `git describe --tag`
COMMIT_HASH = `git rev-parse --short HEAD`
TEST_PATH = `go list ./... | grep -v github.com/QuickAmethyst/monosvc/app`

start: build run

.PHONY: build
build:
	go build -ldflags "\
		-X main.Namespace=${NAMESPACE} \
		-X main.BuildTime=${BUILD_TIME} \
		-X main.BuildVersion=${BUILD_VERSION} \
		-X main.CommitHash=${COMMIT_HASH}" \
		-race -o ./build/app ./app
coverage:
	go test ${TEST_PATH} -covermode=count -coverpkg=./... -coverprofile=coverage.out -failfast -timeout 900s

coverage-visual: coverage
	go tool cover -html=coverage.out

kill-process:
	lsof -i :8080 | awk '/app/ {print $2}' | xargs kill

run: kill-process
	./build/app

graph-generate:
	go run github.com/99designs/gqlgen generate
