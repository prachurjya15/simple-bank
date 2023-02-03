postgres:
	docker run -p 5432:5432 --name postgres15 --network bank-network -e  POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine 
createdb:
	 docker exec -t postgres15 createdb --username=root --owner=root simple_bank
dropdb:
	 docker exec -t postgres15 dropdb simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1
migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go