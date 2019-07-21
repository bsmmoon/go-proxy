init:
	dep init

install:
	dep ensure

build:
	GOOS=linux CGO_ENABLED=0 go install -ldflags="-s -w" ./...

docker-build:
	docker build .

docker-run:
	docker run \
		$(shell docker images --filter "label=project=go-proxy" --filter "label=image=build" --format="{{.ID}}" | head -1)

docker-clean:
	-docker rm  \
		$(shell docker ps -f "status=exited" -q)
	-docker rmi \
		$(shell docker images --filter "label=project=go-proxy" --filter "label=image=build" --format="{{.ID}}" )
