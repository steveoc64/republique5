package republique

import (
	context "context"
	"errors"

	"github.com/sirupsen/logrus"
)

func (s *Server) Login(c context.Context, req *LoginMessage) (*LoginResponse, error) {
	s.log.WithFields(logrus.Fields{
		"Access": req.AccessCode,
		"Team":   req.TeamCode,
		"Player": req.PlayerCode,
	}).Println("Login gRPC")
	// check the game access code
	if req.AccessCode != s.game.AccessCode {
		s.log.Error("Invalid Access Code")
		return nil, errors.New("Invalid Access Code")
	}
	// check the team code
	for _, team := range s.game.Scenario.GetTeams() {
		if team.AccessCode == req.TeamCode {
			// check the player code
			for _, player := range team.Players {
				if player.AccessCode == req.PlayerCode {
					rsp := &LoginResponse{
						Welcome:    "welcome",
						Commanders: player.Commanders,
						TeamName:   team.Name,
						Briefing:   team.Briefing,
					}
					s.log.WithFields(logrus.Fields{
						"Commanders": rsp.Commanders,
						"Team":       rsp.TeamName,
						"Briefing":   rsp.Briefing,
					})
					return rsp, nil
				}
			}
		}
	}
	return nil, errors.New("Invalid Team/Player Code")
}
