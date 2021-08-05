build:
	go build -o battlesnake ./cmd/battlesnake/*.go

docker-build:
	docker build . -t docker-registry.apps.lockleartech.com/clocklear-battlesnake:latest

docker-push:
	docker push docker-registry.apps.lockleartech.com/clocklear-battlesnake:latest

vendor:
	go mod tidy && go mod vendor

.PHONY: build docker-build docker-push vendor