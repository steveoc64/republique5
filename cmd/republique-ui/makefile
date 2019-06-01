all:  install

install: bundle.go
	go install ./...

run: install
	republique-ui

bundle.go: *.jpg
	#fyne bundle republique.png > bundle.go
	fyne bundle icon.jpg > bundle.go

build: bundle.go
	go build ./...

package-darwin: bundle.go
	GOOS=darwin go build ./...
	fyne package -os darwin -icon icon.jpg

package-linux: bundle.go
	GOOS=linux go build ./...
	fyne package -os linux -icon icon.jpg

