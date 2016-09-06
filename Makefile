default: test

# go get github.com/cespare/reflex
.PHONY: watch
watch:
	reflex -r '\.go$$' -- make watch.run

.PHONY: watch.run
watch.run:
	clear
	@make test
	@make lint

.PHONY: test
test:
	go test ./...

# go get -u github.com/golang/lint/golint
.PHONY: lint
lint:
	golint ./...
