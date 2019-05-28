package main

import (
	"fmt"
	"strings"
	"time"

	rp "github.com/steveoc64/republique5/proto"
	"github.com/steveoc64/republique5/republique"
	"github.com/steveoc64/republique5/republique/db"

	"github.com/sirupsen/logrus"
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
	println("Game:", game, "AccessCode =", data.AccessCode)
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
						nn := ""
						switch unit.Arm {
						case rp.Arm_CAVALRY:
							nn = fmt.Sprintf("(%d horse)", unit.Strength*300)
						case rp.Arm_INFANTRY:
							nn = fmt.Sprintf("(%d men)", unit.Strength*550)
							if unit.SkirmisherMax > 0 {
								nn = nn + fmt.Sprintf(" [%d sk]", unit.SkirmisherMax)
							}
						}
						println(unit.Name, strings.ToLower(unit.Grade.String()), strings.Replace(strings.ToLower(unit.UnitType.String()), "_", " ", 1), nn)
					}
					for _, subCommand := range c.Subcommands {
						print("        (", subCommand.Id, ") ", subCommand.Name, " - ", subCommand.CommanderName)
						println("")
						for _, unit := range subCommand.Units {
							print("          - (", unit.Id, ") ")
							nn := ""
							switch unit.Arm {
							case rp.Arm_CAVALRY:
								nn = fmt.Sprintf("(%d horse)", unit.Strength*300)
							case rp.Arm_INFANTRY:
								nn = fmt.Sprintf("(%d men)", unit.Strength*550)
								if unit.SkirmisherMax > 0 {
									nn = nn + fmt.Sprintf(" [%d sk]", unit.SkirmisherMax)
								}
							}
							println(unit.Name, strings.ToLower(unit.Grade.String()), strings.Replace(strings.ToLower(unit.UnitType.String()), "_", " ", 1), nn)
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
