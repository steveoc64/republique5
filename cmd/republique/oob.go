package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/steveoc64/republique5/republique"
)

func oob(log *logrus.Logger, game string) error {
	if !strings.HasSuffix(game, ".db") {
		game = game + ".db"
	}
	db, err := republique.OpenDB(log, game)
	if err != nil {
		return err
	}
	data := &republique.Game{}
	err = db.Load(data)
	if err != nil {
		return err
	}
	println("Game:", game, "AccessCode =", data.AccessCode)
	println("Date:", time.Unix(data.GameTime.Seconds, 0).UTC().Format("Monday, 02-Jan-2006 15:04"))
	print("Table: ", data.TableX, "x", data.TableY, " ft tabletop\n")
	for _, team := range data.Scenario.GetTeams() {
		println("  =============================================================================")
		println("  Team", team.Name, "AccessCode =", team.AccessCode)
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
						nn := ""
						switch unit.Arm {
						case republique.Arm_CAVALRY:
							nn = fmt.Sprintf("(%d horse)", unit.Strength*550)
						case republique.Arm_INFANTRY:
							nn = fmt.Sprintf("(%d men)", unit.Strength*550)
						}
						println(unit.Name, strings.ToLower(unit.Grade.String()), strings.ToLower(unit.UnitType.String()), nn)
					}
					for _, subCommand := range c.Subcommands {
						print("        (", subCommand.Id, ") ", subCommand.Name, " - ", subCommand.CommanderName)
						println("")
						for _, unit := range subCommand.Units {
							print("          - (", unit.Id, ") ")
							nn := ""
							switch unit.Arm {
							case republique.Arm_CAVALRY:
								nn = fmt.Sprintf("(%d horse)", unit.Strength*550)
							case republique.Arm_INFANTRY:
								nn = fmt.Sprintf("(%d men)", unit.Strength*550)
							}
							println(unit.Name, strings.ToLower(unit.Grade.String()), strings.ToLower(unit.UnitType.String()), nn)
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
