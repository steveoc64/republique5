all: build

build:
	go build cmd/rserver/...

run:
	cd cmd/rserver
