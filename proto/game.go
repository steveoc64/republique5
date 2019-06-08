package republique

// GenerateIDs loops through all the commanders and units in a game
// and stamps unique sequential IDs on each unit.
// This should only be used after compilation.
func (g *Game) GenerateIDs() {
	if g.TurnNumber > 0 {
		return
	}
	var cID, ccID, uID int32
	for _, team := range g.GetScenario().GetTeams() {
		for _, command := range team.GetCommands() {
			cID += 100
			command.Id = cID
			ccID = 0
			uID = 1
			for _, unit := range command.Units {
				unit.Id = cID + uID
				uID++
			}
			ccID = 10
			uID = 1
			for _, subCommand := range command.GetSubcommands() {
				subCommand.Id = cID + ccID
				uID = 1
				for _, unit := range subCommand.Units {
					unit.Id = cID + ccID + uID
					uID++
				}
				ccID += 10
			}

		}

	}
}

// InitGameState will loop through the commanders and units in a game
// and set the GameState to defaults based upon the arrival data.
// This is only intended to be run on loading a game in the game server.
// It will return with no changes if the game has already started.
func (g *Game) InitGameState() {
	if g.TurnNumber > 0 {
		return
	}
	for _, team := range g.GetScenario().GetTeams() {
		for _, command := range team.GetCommands() {
			command.initState(nil, false)
			for _, unit := range command.Units {
				unit.initState(command, false)
			}
			numInf := 0
			for _, subCommand := range command.GetSubcommands() {
				standDown := false
				if subCommand.Arm == Arm_INFANTRY {
					numInf++
					if numInf > 1 {
						standDown = true
					}
				}
				subCommand.initState(command, standDown)
				for _, unit := range subCommand.Units {
					unit.initState(command, standDown)
				}
			}

		}

	}
}

// GetUnitCommander returns the command in charge of the given unit
func (g *Game) GetUnitCommander(unit *Unit) *Command {
	for _, team := range g.Scenario.Teams {
		for _, command := range team.Commands {
			for _, u := range command.Units {
				if u.Id == unit.Id {
					return command
				}
			}
			for _, subCommand := range command.Subcommands {
				for _, u := range subCommand.Units {
					if u.Id == unit.Id {
						return subCommand
					}
				}
			}
		}
	}
	return nil
}
