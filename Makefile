password=secret

mysql:
	docker run --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=$(password) -e MYSQL_DATABASE=movie_hub -d mysql:8.4
createdb:
	docker exec -i mysql mysql -u root -p$(password) --execute "create database movie_hub;"
dropdb:
	docker exec -i mysql mysql -u root -p$(password) --execute "drop database movie_hub;"
create_table:
	docker exec -i mysql mysql -u root -p$(password) -D movie_hub < ./movie-hub.sql
sqlc:
	sqlc generate
test:
	go test -v -cover -short ./...
server:
	go run main.go
.PHONY: mysql createdb dropdb create_table sqlc test server