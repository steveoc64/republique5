package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/steveoc64/republique5/db"
	rp "github.com/steveoc64/republique5/proto"
	"github.com/steveoc64/republique5/republique"
)

func info(log *logrus.Logger, game string, full bool, short bool) error {
	if !strings.HasSuffix(game, ".db") {
		game = game + ".db"
	}
	db, err := db.OpenReadDB(log, game)
	if err != nil {
		return err
	}
	defer db.Close()
	data := &rp.Game{}
	err = db.Load("game", "state", data)
	if err != nil {
		return err
	}
	if full {
		println("Game:", game)
		println("Name:", data.Name)
		println("Date:", time.Unix(data.GameTime.Seconds, 0).UTC().Format(republique.DateTimeFormat))
		print("Table: ", data.TableX, "x", data.TableY, " ft tabletop\n")
	} else {
		fmt.Printf("%-16s %-32s %s\n",
			game,
			data.Name,
			time.Unix(data.GameTime.Seconds, 0).UTC().Format(republique.DateFormat),
		)
		if short {
			return nil
		}
	}

	if full {
		println("  -------------------------------------------------------------------------")
		println("  Admin Access =", data.AdminAccess)
	}
	for _, team := range data.Scenario.GetTeams() {
		println("  -------------------------------------------------------------------------")
		if full {
			println("  Team", team.Name, "AccessCode =", team.AccessCode, "GameName =", team.GameName)
			println("")
		} else {
			println("  Team", team.Name, ":", team.AccessCode)
		}
		for _, player := range team.GetPlayers() {
			if full {
				println("    Player AccessCode = ", player.GetAccessCode())
			} else {
				print("    ", player.GetAccessCode(), " (")
			}
			for ii, commander := range player.GetCommanders() {
				c := team.GetCommandByCommanderName(commander)
				if full {
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
				} else {
					if ii > 0 {
						print(", ")
					}
					print(commander)
				}
			}
			if !full {
				print(")")
			}
			println("")
		}
	}
	return nil
}
