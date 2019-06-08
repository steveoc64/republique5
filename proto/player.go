package republique

// GetTeam returns the Team that the given player belongs to
func (g *Game) GetTeam(p *Player) *Team {
	for _, team := range g.GetScenario().GetTeams() {
		for _, player := range team.GetPlayers() {
			if player.GetToken() == p.GetToken() {
				return team
			}
		}
	}
	return nil
}
