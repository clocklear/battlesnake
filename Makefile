docker-build:
	docker build . -t docker-registry.apps.lockleartech.com/clocklear-battlesnake:latest

docker-push:
	docker push docker-registry.apps.lockleartech.com/clocklear-battlesnake:latest

.PHONY: docker-build docker-push