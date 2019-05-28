package republique

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	rp "github.com/steveoc64/republique5/proto"
	"time"
)

func NewToken() *rp.Token {
	return &rp.Token{
		Id:      uuid.New().String(),
		Expires: &timestamp.Timestamp{Seconds: time.Now().UTC().Add(time.Hour * 24).Unix()},
	}
}
