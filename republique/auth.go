package republique

import (
	"github.com/sirupsen/logrus"
	"time"

	"github.com/micro/protobuf/ptypes"
)

func (s *Server) Auth(token string) error {
	for _, team := range s.game.GetScenario().GetTeams() {
		for _, player := range team.GetPlayers() {
			if player.GetToken().GetId() == token {
				t, err := ptypes.Timestamp(player.GetToken().GetExpires())
				if err != nil {
					s.log.WithFields(logrus.Fields{
						"token":   player.GetToken().GetId(),
						"expires": player.GetToken().GetExpires(),
					}).WithError(err).Error("parsing expired value of token")
					return errUnauthorised
				}
				if time.Now().After(t) {
					return errUnauthorised
				}
				return nil
			}
		}
	}
	return errUnauthorised
}
