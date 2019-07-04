all:  install

install:
	go install ./...

arm:
	GOOS=linux GOARM=7 GOARCH=arm go build .
