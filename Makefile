

.PHONY: build
build:
	@go build ./...

.PHONY: format
format:
	@gofumpt -l -w .
	@find . \( -iname '*.c' -o -iname '*.h' \) -type f -print0 | xargs -0 -L1 clang-format -i --style=Google
