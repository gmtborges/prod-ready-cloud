
dev:
	@air

run:
	@go run .

test:
	@go run test -v ./...

o11y-up:
	@cd infra/o11y && docker compose up -d

o11y-down:
	@cd infra/o11y && docker compose down
