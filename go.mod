module github.com/steveoc64/republique5

go 1.12

require (
	fyne.io/fyne v1.1.1
	github.com/boltdb/bolt v1.3.1
	github.com/davecgh/go-spew v1.1.1
	github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d
	github.com/golang/protobuf v1.3.2
	github.com/google/uuid v1.1.1
	github.com/grpc-ecosystem/grpc-gateway v1.9.6
	github.com/hajimehoshi/go-mp3 v0.2.1
	github.com/hajimehoshi/oto v0.4.0
	github.com/llgcode/draw2d v0.0.0-20190810100245-79e59b6b8fbc
	github.com/micro/protobuf v0.0.0-20180321161605-ebd3be6d4fdb
	github.com/sirupsen/logrus v1.4.2
	github.com/steveoc64/memdebug v1.0.0
	github.com/vdobler/chart v1.0.0
	golang.org/x/crypto v0.0.0-20190820162420-60c769a6c586
	google.golang.org/grpc v1.23.0
	k8s.io/apimachinery v0.0.0-20190827074644-f378a67c6af3
)

replace fyne.io/fyne v1.1.1 => ../../../fyne.io/fyne
