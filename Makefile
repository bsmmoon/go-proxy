-include config/.env

init:
	-dep init
	cp -r config/example/. config/

install:
	go get -t -d \
		github.com/gookit/color \
		gopkg.in/elazarl/goproxy.v1
	dep ensure

build:
	GOOS=darwin CGO_ENABLED=0 go install -ldflags="-s -w" ./...

build-window:
	cmd /c /v "set GOOS=windows&& set CGO_ENABLED=0&& go install -ldflags="-s -w" ./..."

run-window:
	cmd /c "set PROXY_PORT=${PROXY_PORT}&${GOPATH}/bin/proxy ${ARGS}"

# Ex. make ARGS="-content-type=image,javascript" run
run:
	PROXY_PORT=${PROXY_PORT} ${GOPATH}/bin/proxy ${ARGS}

clean-output:
	rm -rf ./output
