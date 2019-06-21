package republique

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"

	rp "github.com/steveoc64/republique5/proto"

	"github.com/sirupsen/logrus"
)

// Login is an RPC that validates the given credentials, and returns a struct with the player data if success
func (s *Server) Login(c context.Context, req *rp.LoginMessage) (*rp.LoginResponse, error) {
	s.log.WithFields(logrus.Fields{
		"Hash":   req.Hash,
		"Team":   req.TeamCode,
		"Player": req.PlayerCode,
	}).Info("Login gRPC")
	// check the team code
	s.Lock()
	defer s.Unlock()
	for _, team := range s.game.Scenario.GetTeams() {
		//if team.AccessCode == req.TeamCode {
		// check the player code
		for _, player := range team.Players {
			//if player.AccessCode == req.PlayerCode {
			pwd := []byte(team.AccessCode + player.AccessCode)
			if err := bcrypt.CompareHashAndPassword([]byte(req.Hash), pwd); err == nil {
				player.Token = rp.NewToken()
				s.mTokenCache.Lock()
				s.tokenCache[player.Token.GetId()] = player
				s.mTokenCache.Unlock()
				s.db.Save("game", "state", s.game)
				rsp := &rp.LoginResponse{
					Welcome:    "welcome",
					Commanders: player.Commanders,
					TeamName:   team.Name,
					Briefing:   team.Briefing,
					GameName:   team.GameName,
					GameTime:   s.game.GameTime,
					Token:      player.GetToken(),
				}
				s.log.WithFields(logrus.Fields{
					"Commanders": rsp.Commanders,
					"Team":       rsp.TeamName,
					"GameName":   rsp.GameName,
					"GameTime":   rsp.GameTime.String(),
					"Token":      rsp.Token.Id,
					"Expires":    rsp.Token.Expires.String(),
				}).Info("Player Login")
				return rsp, nil
			}
		}
		//}
	}
	return nil, errors.New("Invalid Team/Player Code")
}
