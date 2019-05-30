package main

import (
	"github.com/steveoc64/republique5/db"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	rp "github.com/steveoc64/republique5/proto"
	"github.com/steveoc64/republique5/republique"
)

func oob(log *logrus.Logger, game string) error {
	if !strings.HasSuffix(game, ".db") {
		game = game + ".db"
	}
	db, err := db.OpenReadDB(log, game)
	if err != nil {
		return err
	}
	data := &rp.Game{}
	err = db.Load("game", "state", data)
	if err != nil {
		return err
	}
	println("Game:", game)
	println("Date:", time.Unix(data.GameTime.Seconds, 0).UTC().Format(republique.DateTimeFormat))
	print("Table: ", data.TableX, "x", data.TableY, " ft tabletop\n")
	for _, team := range data.Scenario.GetTeams() {
		println("  =============================================================================")
		println("  Team", team.Name, "AccessCode =", team.AccessCode, "GameName =", team.GameName)
		println("")
		for _, player := range team.GetPlayers() {
			println("    Player AccessCode =", player.GetAccessCode())
			for _, commander := range player.GetCommanders() {
				c := team.GetCommandByCommanderName(commander)
				if c != nil {
					println("")
					print("      (", c.Id, ") ", c.Name, " - ", commander, ": ")
					println("")
					for _, unit := range c.Units {
						print("      - (", unit.Id, ") ")
						print(unit.LabelString())
					}
					for _, subCommand := range c.Subcommands {
						print("        (", subCommand.Id, ") ", subCommand.Name, " - ", subCommand.CommanderName)
						println("")
						for _, unit := range subCommand.Units {
							print("          - (", unit.Id, ") ")
							println(unit.LabelString())
						}
					}
				}
			}
			println("")
			println("  -------------------------------------------------------------------------")
		}
	}
	return nil
}
