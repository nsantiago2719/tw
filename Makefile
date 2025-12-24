# Go parameters
GOCMD = go
GOTEST = $(GOCMD) test

# Coverage
COVERAGE_REPORT = coverage.out
COVERAGE_MODE = count

.PHONY: build test run release test-coverage

build:
	@go build -o bin/tw ./cmd/main.go

test:
	go test ./...

run: build
	./bin/tw

release:
	CGO_ENABLED=0 go build -ldflags "-s -w" -o /release/tw ./cmd/main.go

test-coverage:
	echo "" > $(COVERAGE_REPORT); \
	$(GOTEST) -coverprofile=$(COVERAGE_REPORT) -coverpkg=./... -covermode=$(COVERAGE_MODE) ./...
