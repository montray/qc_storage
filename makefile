#!make
include .env

start:
	@docker-compose up -d
	@echo dump sql...
	@(docker exec -i $(PG_CONTAINER_NAME) sh -c "cat > dump.sql") < dump.sql
	@docker exec -i $(PG_CONTAINER_NAME) sh -c "psql storage admin < dump.sql"
	@echo building app...
	@rm -f out/server
	@go build -o out/server cmd/main.go
	@echo starting server...
	@out/server

stop:
	@docker-compose stop
