#!make

start:
	@docker-compose up -d
	@echo dump sql...
	@docker exec -it storage_pg /bin/bash -c echo dump.sql
	@echo building app...
	@rm -f out/server
	@go build -o out/server cmd/main.go
	@echo starting server...
	@out/server

stop:
	@docker-compose stop
