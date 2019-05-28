package appwindow

import (
	"context"
	rp "github.com/steveoc64/republique5/proto"
	"time"
)

func (a *App) Phaser() {
	go func() {
		for {
			time.Sleep(time.Second * 3)

			_, err := a.client.GameTime(context.Background(), &rp.StringMessage{Value: a.Token})
			if err != nil {
				println("got an error", err.Error())
			}
			// TODO - if any changes then write them to the screen, and save the data into the app struct

		}
	}()
}
