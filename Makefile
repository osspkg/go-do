
.PHONY: install
install:
	go install github.com/osspkg/devtool@latest

.PHONY: setup
setup:
	devtool setup-lib

.PHONY: lint
lint:
	devtool lint

.PHONY: license
license:
	devtool license

.PHONY: tests
tests:
	devtool test

.PHONY: pre-commit
pre-commit: setup license lint tests

.PHONY: ci
ci: install setup lint tests

