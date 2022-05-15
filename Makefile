test: lint unit-test

lint:
	@docker build --target lint .

unit-test:
	@docker build --target coverage-test --output ./ .
