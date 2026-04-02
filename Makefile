include .env
export

export PROJECT_ROOT=${shell pwd}

env-up:
	@docker compose up -d todoapp-postgres
env-down:
	@docker compose down todoapp-postgres
env-cleanup:
	@read -p "Очистити всі volume файли? [y/N]" ans; \
	if [ "$$ans" = "y" ]; then \
	  docker compose down todoapp-postgres port-forwarder && \
	  sudo rm -rf ${PROJECT_ROOT}/out/pgdata && \
	  echo "Файли очищені"; \
	else \
	  echo "Очистка середовища скасована"; \
	fi

migrate-create:
	@if [ -z "$(seq)" ]; then \
  		echo "Відсутній seq"; \
  		exit 1; \
  	fi;
	docker compose run --rm todoapp-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
  		echo "Відстуній action"; \
  		exit 1; \
  	fi;
	docker compose run --rm todoapp-postgres-migrate \
    	-path /migrations \
    	-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=disable \
    	"$(action)"

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose dowm port-forwarder

todoapp-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	go run ${PROJECT_ROOT}/cmd/todoapp/main.go