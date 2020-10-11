default:
	@echo "=============building Local API============="
	docker build -f Dockerfile -t portfolio-api .

up: default
	@echo "=============starting api locally============="
	docker-compose up -d

windows-up:
	@echo "=============starting api locally, no docker============="
	CompileDaemon -build="go build -o portfolio-server.exe ./src" -command="./portfolio-server"

logs:
	docker-compose logs -f

down:
	docker-compose down

test:
	go test -v -cover ./...

clean: down
	@echo "=============cleaning up============="
	rm main
	docker system prune -f
	docker volume prune -f
