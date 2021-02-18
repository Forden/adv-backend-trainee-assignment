.PHONY: all clean restart_api migrate

all:
	docker run -d --name avito-db --network host -e POSTGRES_USER=username -e POSTGRES_PASSWORD=password -e POSTGRES_DB=avito-backend -e PGDATA=./db/postgres postgres:13
	docker build -t avito/adv-api .
	docker run -d --name avito-api --network host -e CONFIG_PATH=/config/config.json -v $(shell pwd)/config/config.json:/config/config.json avito/adv-api

clean:
	docker stop avito-db || true && docker rm avito-db || true
	docker stop avito-api || true && docker rm avito-api || true

restart_api:
	docker stop avito-api || true && docker rm avito-api || true
	docker build -t avito/adv-api .
	docker run -d --name avito-api --network host -e CONFIG_PATH=/config/config.json -v $(shell pwd)/config/config.json:/config/config.json avito/adv-api

migrate:
	docker run -v $(shell pwd)/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database postgres://username:password@localhost:5432/avito-backend?sslmode=disable up