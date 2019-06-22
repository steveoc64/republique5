package republique

// GetCommandByName gets the matching command by the name of the unit
func (t *Team) GetCommandByName(name string) *Command {
	for _, command := range t.GetCommands() {
		if command.Name == name {
			return command
		}
		for _, subCommand := range command.Subcommands {
			if subCommand.Name == name {
				return subCommand
			}
		}
	}
	return nil
}

// GetCommandByID gets the matching command by the ID of the unit
func (t *Team) GetCommandByID(id int32) *Command {
	for _, command := range t.GetCommands() {
		if command.Id == id {
			return command
		}
		for _, subCommand := range command.Subcommands {
			if subCommand.Id == id {
				return subCommand
			}
		}
	}
	return nil
}

// GetCommandByCommanderName gets the matching command by the name of the commander
func (t *Team) GetCommandByCommanderName(name string) *Command {
	for _, command := range t.GetCommands() {
		if command.CommanderName == name {
			return command
		}
		for _, subCommand := range command.Subcommands {
			if subCommand.CommanderName == "" {
				continue
			}
			if subCommand.CommanderName == name {
				return subCommand
			}
		}
	}
	return nil
}
