BINARY_NAME=rate-limiter
MAKEFLAGS += --silent

PROD_ENV = .env.prod
DEV_ENV = .env.dev
ifeq ($(MAKECMDGOALS),run-prod)
	include $(PROD_ENV)
else
	include $(DEV_ENV)
endif

APP_ARGUMENTS = -db_user ${DB_USER} \
                -db_pass ${DB_PASS}	\
                -server_host ${SERVER_HOST} \
                -server_port ${SERVER_PORT} \
                -db_host ${DB_HOST} \
                -db_port ${DB_PORT} \
                -db_name ${DB_NAME} \
                -db_user ${DB_USER} \
                -db_pass ${DB_PASS}

run: db-up run-dev

build:
	echo "Building..."
	GOARCH=amd64 go build -o ${BINARY_NAME} .

db-up:
	echo "Starting database..."
	docker compose -f docker/docker-compose.dev-db.yml up -d

db-down:
	echo "Stopping database..."
	docker compose -f docker/docker-compose.dev-db.yml down

run-dev: build
	echo "Launching DEV app..."
	./${BINARY_NAME} $(APP_ARGUMENTS)

run-prod: build
	echo "Launching PROD app..."
	./${BINARY_NAME} $(APP_ARGUMENTS)

clean: db-down
	go clean
	rm ${BINARY_NAME}

test:
	go test ./...

dep:
	go mod download

lint:
	golangci-lint run
