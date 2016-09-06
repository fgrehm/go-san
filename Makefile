default: test

# go get github.com/cespare/reflex
.PHONY: watch
watch:
	reflex -r '\.go$$' -- make watch.run

.PHONY: watch.run
watch.run:
	clear
	@$(MAKE) test

.PHONY: test
test:
	go test ./...
