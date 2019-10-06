-include config/.env

init:
	-dep init
	cp -r config/example/. config/

install:
	go get -t -d \
		github.com/tebeka/selenium \
		github.com/gookit/color \
		gopkg.in/elazarl/goproxy.v1
	dep ensure

build:
	GOOS=darwin CGO_ENABLED=0 go install -ldflags="-s -w" ./...

# Ex. make ARGS="-content-type=image,javascript" run
run:
	SELENIUM_PATH=${SELENIUM_PATH} \
	GECKO_PATH=${GECKO_PATH} \
	CHROME_PATH=${CHROME_PATH} \
	SELENIUM_PORT=${SELENIUM_PORT} \
	PROXY_PORT=${PROXY_PORT} \
		${GOPATH}/bin/proxy ${ARGS}

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

clean-output:
	rm -rf ./output
