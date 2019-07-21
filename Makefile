-include config/.env

init:
	-dep init
	cp -r config/example/. config/

install:
	go get -t -d \
		github.com/tebeka/selenium \
		github.com/gookit/color
	dep ensure

build:
	GOOS=darwin CGO_ENABLED=0 go install -ldflags="-s -w" ./...

run:
	SELENIUM_PATH=${SELENIUM_PATH} \
	GECKO_PATH=${GECKO_PATH} \
	PORT=${PORT} \
		${GOPATH}/bin/proxy

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

drivers-update:
	( cd ${GOPATH}/src/github.com/tebeka/selenium/vendor; go get -d .; go run init.go --alsologtostderr )

drivers-list:
	ls -al ${DRIVER_PATH}
