postgres:
	docker run --name postgres_new -p 9876:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

sqlc:
	sqlc generate

migrateup:
	@source .env && migrate -path db/migrations -database "$${DB_SOURCE}" -verbose up

migratedown:
	@source .env && migrate -path db/migrations -database "$${DB_SOURCE}" -verbose down



.PHONY : postgres