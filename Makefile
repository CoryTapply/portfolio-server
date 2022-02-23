default:
	@echo "============= Building Local API ============="
	docker build -f dev.Dockerfile -t portfolio-api-dev .
	docker build -f ./database/Dockerfile -t portfolio-db .

up: default
	@echo "============= Starting API Locally ============="
	docker-compose up -d api-dev
	docker-compose up -d portfolio-db

build-prod:
	@echo "============= Building PRODUCTION API ============="
	docker build -f Dockerfile -t gcr.io/poetic-bison-300502/portfolio-api:v1 .

up-prod: build-prod
	@echo "============= Starting PRODUCTION API Locally ============="
	docker-compose up -d api

windows-up:
	@echo "============= Starting API Locally, No Docker ============="
	CompileDaemon -build="go build -o portfolio-server.exe ./src" -command="./portfolio-server"

logs:
	docker-compose logs -f

down:
	docker-compose down

test:
	go test -v -cover ./...

clean: down
	@echo "============= Cleaning Up ============="
	rm portfolio-server || true
	docker system prune -f
	docker volume prune -f

int:
	@echo "============= Running API with Interactive Shell ============="
	docker run -it portfolio-api --entrypoint sh

i:
	@echo "============= Running API with Interactive Shell ============="
	docker run -it portfolio-api sh

db-i:
	@echo "============= Running DB with Interactive Shell ============="
	# docker exec -it portfolio-db psql -U postgres
	# docker exec -it portfolio-server_portfolio-db_1  psql -U postgres
	# docker exec -it portfolio-server_portfolio-db_1  bash
	docker-compose run portfolio-db bash
