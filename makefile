all: rserver

rserver:
	cd cmd/rserver && go build ./...

run:
	cd cmd/rserver
