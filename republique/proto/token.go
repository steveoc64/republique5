package republique

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	"time"
)

func NewToken() *Token {
	return &Token{
		Id:      uuid.New().String(),
		Expires: &timestamp.Timestamp{Seconds: time.Now().UTC().Add(time.Hour * 24).Unix()},
	}
}
