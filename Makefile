build:
	go build

build-gateway-os:
	GOOS=linux GOARCH=arm CGO=0 go build
	cd packaging/gateway-os && ./package.sh 1.0.0
