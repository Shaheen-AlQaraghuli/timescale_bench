make setup:
	docker-compose up -d

migrate-up:
	goose -dir=./database/timescale/migrations -allow-missing postgres "host=localhost user=postgres password=postgres dbname=shift_db port=5432 sslmode=disable TimeZone=UTC" up
	goose -dir=./database/postgres/migrations -allow-missing postgres "host=localhost user=postgres password=postgres dbname=shift_db port=5433 sslmode=disable TimeZone=UTC" up

benchmark:
	go test -bench=. ./main