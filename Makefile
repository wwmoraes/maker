# Disable built-in rules and variables and suffixes
MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --no-builtin-variables
.SUFFIXES:

# chains common rule names to included ones
.DEFAULT_GOAL := all
.PHONY: all build clean test coverage lint run
all: build lint test coverage
build: golang-build
clean: golang-clean
test: golang-test
coverage: golang-coverage golang-coverage-html
lint: golang-lint
diagrams: goplantuml-run
docs: plantuml-render

plantuml-render: goplantuml-run

run: GOLANG_RUN:=./cmd/maker
run: GOLANG_RUN_ARGS:=$(ARGS)
run: golang-run

cli-add: GOLANG_RUN:=./cmd/maker
cli-add: GOLANG_RUN_ARGS:=add $(ARGS)
cli-add: golang-run

cli-remove: GOLANG_RUN:=./cmd/maker
cli-remove: GOLANG_RUN_ARGS:=remove $(ARGS)
cli-remove: golang-run

cli-install: GOLANG_RUN:=./cmd/maker
cli-install: GOLANG_RUN_ARGS:=install $(ARGS)
cli-install: golang-run

ctl-repo-init: GOLANG_RUN:=./cmd/makerctl
ctl-repo-init: GOLANG_RUN_ARGS:=repository init $(ARGS)
ctl-repo-init: golang-run

ctl-repo-index: GOLANG_RUN:=./cmd/makerctl
ctl-repo-index: GOLANG_RUN_ARGS:=repository index
ctl-repo-index: golang-run

ctl-snippet-init: GOLANG_RUN:=./cmd/makerctl
ctl-snippet-init: GOLANG_RUN_ARGS:=snippet init $(ARGS)
ctl-snippet-init: golang-run

ctl-snippet-release: GOLANG_RUN:=./cmd/makerctl
ctl-snippet-release: GOLANG_RUN_ARGS:=snippet release $(ARGS)
ctl-snippet-release: golang-run

################################################################################
### variables and includes
################################################################################

# pre-include variables
DEBUG := 1
GO := grc go

# includes
include .make/*.mk

# post-include variables
# override B_VAR += value
