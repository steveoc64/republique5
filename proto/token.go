package republique

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	"time"
)

// NewToken returns a new token with a randomised UUID and a 24hr expiry period
func NewToken() *Token {
	return &Token{
		Id:      uuid.New().String(),
		Expires: &timestamp.Timestamp{Seconds: time.Now().UTC().Add(time.Hour * 24).Unix()},
	}
}
