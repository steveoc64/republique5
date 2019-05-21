package main

import (
	"github.com/sirupsen/logrus"
	"github.com/steveoc64/republique5/republique"
	"strings"
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
	println("  Admin Access =", data.AdminAccess)
	for _, team := range data.Scenario.GetTeams() {
		println("  Team", team.Name, "AccessCode =", team.AccessCode)
		for _, player := range team.GetPlayers() {
			println("    Player AccessCode =", player.GetAccessCode())
			for _, commander := range player.GetCommanders() {
				println("      -", commander)
			}

		}
	}
	return nil
}
