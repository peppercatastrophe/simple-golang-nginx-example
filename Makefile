IMAGE_PREFIX := simple-golang-nginx-example

.PHONY: build up dev dev-down rebuild prune clean

build:
	podman build --target runtime -t $(IMAGE_PREFIX)-main-service ./main-service
	podman build --target runtime -t $(IMAGE_PREFIX)-items-service ./items-service
	podman image prune -f

up:
	podman compose up

dev:
	podman compose -f docker-compose.yml -f docker-compose.dev.yml up

dev-down:
	podman compose -f docker-compose.yml -f docker-compose.dev.yml down

rebuild: build up

prune:
	podman image prune -f

clean:
	podman compose down -v 2>/dev/null; \
	podman image prune -f; \
	podman image rm $(IMAGE_PREFIX)-main-service $(IMAGE_PREFIX)-items-service 2>/dev/null; \
	podman image rm golang:1.26-alpine 2>/dev/null; \
	echo "Cleaned up"
