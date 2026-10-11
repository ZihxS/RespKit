.PHONY: test race vet coverage integration integration-mappers doctor

test:
	go test -count=1 ./...

race:
	go test -race -count=1 ./...

vet:
	go vet ./...

coverage:
	go test -count=1 ./resp ./internal/... -coverpkg=./resp,./internal/... -coverprofile=coverage.out
	go tool cover -func=coverage.out
	@go tool cover -func=coverage.out | grep -E '^total:.*100\.0%$$'

integration:
	cd integrations/framework-tests && go mod tidy && go test -race -count=1 -v ./...

integration-mappers:
	cd integrations/pgx && go test -race -count=1 ./...
	cd integrations/sqlite3 && go test -race -count=1 -cover ./... && go vet ./...

doctor:
	go run ./cmd/respkit doctor
