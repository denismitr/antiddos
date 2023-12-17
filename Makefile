.PHONY: server/run
server/run:
	go run cmd/server/main.go

.PHONY: client/run
client/run:
	go run cmd/client/main.go

.PHONY: test
test:
	go test ./...

.PHONY: docker
docker:
	docker-compose -f zarf/docker/docker-compose.yml up -d --build

.PHONY: docker/clean
docker/clean:
	docker-compose -f zarf/docker/docker-compose.yml rm --force --stop -v