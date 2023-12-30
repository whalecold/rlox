# Golang standard bin directory.
GOPATH ?= $(shell go env GOPATH)
BIN_DIR := $(GOPATH)/bin
GOIMPORTS=$(BIN_DIR)/goimports

.PHONY: tool build

.PHONY: format
format: ast
ifeq ("$(wildcard $(GOIMPORTS))", "")
	go install golang.org/x/tools/cmd/goimports@latest
endif
	@$(GOIMPORTS) -w ./

ast: 
	go run tool/tool.go ./

build:
	go build -o glox  *.go

install:
	go install  *.go


