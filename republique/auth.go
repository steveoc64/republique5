package republique

import (
	"github.com/davecgh/go-spew/spew"
	"time"

	"github.com/micro/protobuf/ptypes"
	"github.com/sirupsen/logrus"
	rp "github.com/steveoc64/republique5/proto"
)

func (s *Server) Auth(token string) (*rp.Player, error) {
	p, ok := s.tokenCache[token]
	if !ok {
		return nil, errUnauthorised
	}
	e, err := ptypes.Timestamp(p.Token.GetExpires())
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"token":   p.Token.GetId(),
			"expires": p.Token.GetExpires(),
		}).WithError(err).Error("parsing expired value of token")
		return nil, errUnauthorised
	}
	if time.Now().After(e) {
		return nil, errSessionExpired
	}
	return p, nil
}

func NewTokenCache(game *rp.Game) map[string]*rp.Player {
	tokenCache := make(map[string]*rp.Player)
	for _, team := range game.GetScenario().GetTeams() {
		for _, player := range team.GetPlayers() {
			if player.Token != nil {
				tokenCache[player.Token.GetId()] = player
			}
		}
	}
	spew.Dump("token cache", tokenCache)
	return tokenCache
}
