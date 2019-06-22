package appwindow

import (
	"context"
	"time"
)

// Phaser is a goroutine that loops on a 3 second timer to fetch the game time
func (a *App) Phaser() {
	go func() {
		for {
			time.Sleep(time.Second * 3)

			_, err := a.gameServer.GameTime(context.Background(), &a.Token)
			if err != nil {
				println("got an error", err.Error())
			}
			// TODO - if any changes then write them to the screen, and save the data into the app struct

		}
	}()
}