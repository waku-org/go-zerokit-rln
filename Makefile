.PHONY: rlnlib

SHELL := bash # the shell used internally by Make

GOBIN ?= $(shell which go)

rlnlib:
	scripts/build.sh
	cd zerokit/rln && cbindgen --config ../../cbindgen.toml --crate rln --output ../../rln/librln.h --lang c

test:
	go test ./... -count 1 -v