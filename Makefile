build:
	go build

build-gateway-os:
	docker compose run --rm chirpstack-pg-to-sqlite-armv7 go build -ldflags "-linkmode 'external' -extldflags '-static'"
	cd packaging/gateway-os && ./package.sh 1.0.0
