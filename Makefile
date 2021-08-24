postgres:
	docker run --name some-postgres -p5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=mysecretpassword -d postgres:alpine

createdb:
	docker exec -it some-postgres createdb --username=postgres --owner=postgres development

dropdb:
	docker exec -it some-postgres dropdb development

test:
	go test -v -cover ./...

migrateup:
	migrate -path ./migrations -database "postgresql://postgres:mysecretpassword@localhost:5432/development?sslmode=disable" -verbose up

migrateup1:
	migrate -path ./migration -database "postgresql://postgres:mysecretpassword@localhost:5432/development?sslmode=disable" -verbose up 1

migratedown:
	migrate -path ./migration -database "postgresql://postgres:mysecretpassword@localhost:5432/development?sslmode=disable" -verbose down

migratedown1:
	migrate -path ./migration -database "postgresql://postgres:mysecretpassword@localhost:5432/development?sslmode=disable" -verbose down 1

server:
	go run ./cmd/api/main.go -cors-trusted-origins="http://localhost:3000"

.PHONY: postgres createdb migrateup migratedown migrateup1 migratedown1 dropdb test server