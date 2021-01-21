default:
	@echo "============= Building Local API ============="
	docker build -f dev.Dockerfile -t portfolio-api-dev .

up: default
	@echo "============= Starting API Locally ============="
	docker-compose up -d api-dev

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
