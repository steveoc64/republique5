package republique

import (
	"github.com/davecgh/go-spew/spew"
	"time"

	"github.com/micro/protobuf/ptypes"
	"github.com/sirupsen/logrus"
	rp "github.com/steveoc64/republique5/proto"
)

func (s *Server) Auth(token string) error {
	t, ok := s.tokenCache[token]
	if !ok {
		return errUnauthorised
	}
	e, err := ptypes.Timestamp(t.GetExpires())
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"token":   t.GetId(),
			"expires": t.GetExpires(),
		}).WithError(err).Error("parsing expired value of token")
		return errUnauthorised
	}
	if time.Now().After(e) {
		return errSessionExpired
	}
	return nil
}

func NewTokenCache(game *rp.Game) map[string]*rp.Token {
	tokenCache := make(map[string]*rp.Token)
	for _, team := range game.GetScenario().GetTeams() {
		for _, player := range team.GetPlayers() {
			if player.Token != nil {
				tokenCache[player.Token.GetId()] = player.Token
			}
		}
	}
	spew.Dump("token cache", tokenCache)
	return tokenCache
}
