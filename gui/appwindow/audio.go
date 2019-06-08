package appwindow

import (
	"io"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

var isPlaying bool

func (a *App) PlayAudio(arm string) {
	return
	/*
		dirname := filepath.Join(os.Getenv("HOME"), "republique", arm)
		files := []string{}
		filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
				return err
			}
			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		})
		if len(files) < 1 {
			return
		}
		i := rand.Intn(len(files))
		audioFile := files[i]
		a.playAudio(audioFile)

	*/
}

func (a *App) PlaySystemAudio(name string) {
	a.playAudio(filepath.Join(os.Getenv("HOME"), "republique", "system", name+".mp3"))
}

func (a *App) playAudio(audioFile string) {
	if isPlaying {
		return
	}
	go func() {
		f, err := os.Open(audioFile)
		if err != nil {
			println("audio error:", err.Error())
			return
		}
		defer f.Close()

		d, err := mp3.NewDecoder(f)
		if err != nil {
			println("audio error:", err.Error())
			return
		}

		if a.audioPort == nil {
			p, err := oto.NewPlayer(d.SampleRate(), 2, 2, 8192)
			if err != nil {
				println("audio error:", err.Error())
				return
			}
			a.audioPort = p
		}

		isPlaying = true
		if _, err := io.Copy(a.audioPort, d); err != nil {
			println("audio error:", err.Error())
		}
		isPlaying = false
	}()
}
