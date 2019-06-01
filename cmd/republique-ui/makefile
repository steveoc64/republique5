all:  install

install: bundle.go
	go install ./...

run: install
	republique-ui

bundle.go: republique.png
	fynce bundle republique.png > bundle.go

build: bundle.go
	go build ./...

package-darwin: bundle.go
	GOOS=darwin go build ./...
	fyne package -os darwin -icon republique.png

package-linux: bundle.go
	GOOS=linux go build ./...
	fyne package -os linux -icon republique.png

