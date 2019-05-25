package republique

import (
	"fmt"
	rp "github.com/steveoc64/republique5/republique/proto"
)

type Session struct {
	ServerName   string
	GameName     string
	LoginDetails *rp.LoginResponse
}

func (s *Session) String() string {
	return fmt.Sprintf("%s: %s\n%s\nCommanders: %v\nTeam: %s\nBriefing: %s",
		s.LoginDetails.GetWelcome(),
		s.GameName,
		s.LoginDetails.GetGameName(),
		s.LoginDetails.GetCommanders(),
		s.LoginDetails.GetTeamName(),
		s.LoginDetails.GetBriefing(),
	)
}
