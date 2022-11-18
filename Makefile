
swagger:
	swag init -g ./cmd/main/main.go -o ./docs

run:
	docker-compose up -d --build --force-recreate

stop:
	docker-compose down