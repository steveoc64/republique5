package republique

import (
	"fmt"
	rp "github.com/steveoc64/republique5/republique/proto"
	"time"
)

type Session struct {
	ServerName   string
	LoginDetails *rp.LoginResponse
	GameName     string
	GameTime     time.Time
}

func (s *Session) String() string {
	return fmt.Sprintf("Game: %s\nCommanders: %v\nTeam: %s\nBriefing: %s",
		s.LoginDetails.GetGameName(),
		s.LoginDetails.GetCommanders(),
		s.LoginDetails.GetTeamName(),
		s.LoginDetails.GetBriefing(),
	)
}
