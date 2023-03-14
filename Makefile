.PHONY: test
test:
		go test -race -v ./...
.PHONY: start
start:
	@make clean
	@make build
	@echo "running main program..."
	@./etsim $(ARGS)
.PHONY: build
build:
	@make test
	@golangci-lint run
	@go build ./etsimcmd/etsim
.PHONY: clean
clean:
	@go clean -i
.PHONY: cover
cover:
	@mkdir .coverage || echo "hidden coverage folder exists"
	@go test -v -cover ./... -coverprofile .coverage/coverage.out
	@go tool cover -html=.coverage/coverage.out -o .coverage/coverage.html
.PHONY: covero
covero:
	@make cover
	@open .coverage/coverage.html
