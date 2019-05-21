package republique

func (g *Game) GenerateIDs() {
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
				ccID += 10
				for _, unit := range subCommand.Units {
					unit.Id = cID + ccID + uID
					uID++
				}
			}

		}

	}
}
