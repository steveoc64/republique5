package appwindow

import (
	"context"
	"github.com/steveoc64/memdebug"
	"github.com/steveoc64/republique5/republique"
	rp "github.com/steveoc64/republique5/republique/proto"
	"time"
)

func (a *App) Phaser() {
	go func() {
		for {
			time.Sleep(time.Second * 3)

			t1 := time.Now()
			gameTime, err := a.client.GameTime(context.Background(), &rp.StringMessage{Value: a.Token})
			if err != nil {
				println("got an error", err.Error())
			}
			// TODO - if any changes then write them to the screen, and save the data into the app struct

		}
	}()
}
