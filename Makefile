start:
	go run cmd/main/main.go

migrate_create:
	~/go/bin/migrate create -ext sql -dir ./migrations -format "20060102150405" $(name)

migrate_up:
	~/go/bin/migrate \
		-path ./migrations/ \
		-database "sqlite3://myapp.db" \
		-verbose up

# usage:
# 	make migrate_down N=1 - for one migration down
migrate_down:
	~/go/bin/migrate \
		-path ./migrations/ \
		-database "sqlite3://myapp.db" \
		-verbose down $(N)