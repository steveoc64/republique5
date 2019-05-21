package republique

func (t *Team) GetCommandByCommanderName(name string) *Command {
	println("checking team", t.Name, "for <", name, ">")
	for _, command := range t.GetCommands() {
		println(command.Name, command.CommanderName)
		if command.CommanderName == name {
			return command
		}
		for _, subCommand := range command.Subcommands {
			if subCommand.CommanderName == "" {
				continue
			}
			println("  ", subCommand.Name, subCommand.CommanderName)
			if subCommand.CommanderName == name {
				return subCommand
			}
		}
	}
	return nil
}
