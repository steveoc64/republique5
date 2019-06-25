package republique

import (
	"context"
	"time"

	"github.com/steveoc64/memdebug"

	rp "github.com/steveoc64/republique5/proto"
)

// GetUnits returns all the units owned by the player with the given token
func (s *Server) GetUnits(c context.Context, req *rp.TokenMessage) (*rp.Units, error) {
	s.RLock()
	defer s.RUnlock()
	player, err := s.Auth(req.Id)
	if err != nil {
		return nil, err
	}
	t1 := time.Now()
	team := s.game.GetTeam(player)
	s.log.WithField("token", req.Id).Println("GetUnits gRPC")
	commands := &rp.Units{}
	for _, commander := range player.GetCommanders() {
		c := team.GetCommandByCommanderName(commander)
		if c != nil {
			if c.Arrival.ComputedTurn <= s.game.TurnNumber {
				commands.Commands = append(commands.Commands, c)
			}
		}
	}
	memdebug.Print(t1, "Got units")
	return commands, nil
}

// GetEnemy returns the visible enemy units
func (s *Server) GetEnemy(c context.Context, req *rp.TokenMessage) (*rp.Units, error) {
	s.RLock()
	defer s.RUnlock()
	player, err := s.Auth(req.Id)
	if err != nil {
		return nil, err
	}
	t1 := time.Now()
	team := s.game.GetTeam(player)
	s.log.WithField("token", req.Id).Println("GetEnemy gRPC")
	units := &rp.Units{}
	teams := s.game.Scenario.GetTeams()

	for _, otherteam := range teams {
		if otherteam != team {
			for _, command := range otherteam.GetCommands() {
				// this command is on the other team then
				println("arrives", command.Arrival.ComputedTurn, "now is", s.game.TurnNumber)
				if command.Arrival.ComputedTurn <= s.game.TurnNumber {
					enemyCommand := &rp.Command{
						Id:            command.Id,
						Name:          command.Name,
						CommanderName: command.CommanderName,
						Rank:          command.Rank,
						Arm:           command.Arm,
						Nationality:   command.Nationality,
						GameState: &rp.CommandGameState{
							Grid: &rp.Grid{
								X: command.GameState.Grid.X,
								Y: command.GameState.Grid.Y,
							},
							Position:  command.GameState.Position,
							Formation: command.GameState.Formation,
						},
					}
					// attach subcommands
					for _, subCommand := range command.Subcommands {
						if subCommand.Arrival.ComputedTurn <= s.game.TurnNumber {
							enemySub := &rp.Command{
								Id:            subCommand.Id,
								Name:          subCommand.Name,
								CommanderName: subCommand.CommanderName,
								Rank:          subCommand.Rank,
								Arm:           subCommand.Arm,
								Nationality:   subCommand.Nationality,
								GameState: &rp.CommandGameState{
									Grid: &rp.Grid{
										X: subCommand.GameState.Grid.X,
										Y: subCommand.GameState.Grid.Y,
									},
									Position:  subCommand.GameState.Position,
									Formation: subCommand.GameState.Formation,
								},
							}
							enemyCommand.Subcommands = append(enemyCommand.Subcommands, enemySub)
						}
					}
					units.Commands = append(units.Commands, enemyCommand)
				}
			}
		}
	}
	memdebug.Print(t1, "Got enemy units")
	return units, nil
}
