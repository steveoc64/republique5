package republique

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
