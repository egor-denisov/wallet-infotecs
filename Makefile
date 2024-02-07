build:
	docker compose up -d db --force-recreate
	docker-compose up --build --force-recreate

run:
	go run cmd/app/main.go

test:
	go test ./...

swag:
	swag init -dir internal/controller/http/v1/ -generalInfo router.go --parseDependency internal/entity/ 