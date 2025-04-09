# Go parameters
GOCMD = go
GOTEST = $(GOCMD) test 

# Coverage
COVERAGE_REPORT = coverage.out
COVERAGE_MODE = count

build:
	@go build -o bin/tw

run: build
	./bin/tw

release:
	CGO_ENABLED=0 go build -ldflags "-s -w" -o /release/tw

test-coverage:
	@echo "running against `git version`"; \
	echo "" > $(COVERAGE_REPORT); \
	$(GOTEST) -coverprofile=$(COVERAGE_REPORT) -coverpkg=./... -covermode=$(COVERAGE_MODE) ./...
