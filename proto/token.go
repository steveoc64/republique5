package republique

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	"time"
)

var TokenMap = map[string]*Token

func NewToken() *Token {
	return &Token{
		Id:      uuid.New().String(),
		Expires: &timestamp.Timestamp{Seconds: time.Now().UTC().Add(time.Hour * 24).Unix()},
	}
}

