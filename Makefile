include .env

.PHONY: lint unit-test build

all: build

test: lint unit-test

lint:
	@docker build --target lint .

unit-test:
	@docker build --target coverage-test --output ./ .

build:
	@docker build --target bin --output bin/ .
