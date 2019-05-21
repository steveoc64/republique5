package main

import (
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/steveoc64/republique5/republique"
)

func info(log *logrus.Logger, game string) error {
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
	println("  -------------------------------------------------------------------------")
	println("  Admin Access =", data.AdminAccess)
	for _, team := range data.Scenario.GetTeams() {
		println("  -------------------------------------------------------------------------")
		println("  Team", team.Name, "AccessCode =", team.AccessCode)
		println("")
		for _, player := range team.GetPlayers() {
			println("    Player AccessCode =", player.GetAccessCode())
			for _, commander := range player.GetCommanders() {
				c := team.GetCommandByCommanderName(commander)
				if c != nil {
					print("      -", commander, ": ")
					print(c.Arrival.Position.String())
					if c.Arrival.From > 0 {
						print(" Arrives ", c.Arrival.From, "-", c.Arrival.To, " hrs")
					}
					if c.Arrival.Percent < 100 {
						print(" with a ", c.Arrival.Percent, "% chance")
					}
					if c.Arrival.Contact {
						print(" ***Contact***")
					}
					println("")
				}
			}
			println("")
		}
	}
	return nil
}
