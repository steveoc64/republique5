all:  install

install:
	go install ./...

mobile:
	#CGO_CFLAGS="-I/home/steve/Android/Sdk/ndk-bundle/sysroot/usr/include -I/usr/include" gomobile build -target android/arm7a .
	gomobile build -target android .

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

