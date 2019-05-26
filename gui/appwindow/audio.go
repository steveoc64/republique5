package appwindow

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

func (a *App) PlayAudio(arm string) {
	dirname := filepath.Join(os.Getenv("HOME"), "republique", arm)
	os.Chdir(dirname)
	files := []string{}
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() && info.Name() == arm {
			return filepath.SkipDir
		}
		files = append(files, path)
		return nil
	})
	if len(files) < 1 {
		return
	}
	i := rand.Intn(len(files))
	audioFile := files[i]
	go a.playAudio(audioFile)
}

func (a *App) playAudio(audioFile string) error {
	f, err := os.Open(audioFile)
	if err != nil {
		return err
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		return err
	}

	if a.player == nil {
		p, err := oto.NewPlayer(d.SampleRate(), 2, 2, 8192)
		if err != nil {
			return err
		}
		a.player = p
	}

	fmt.Printf("Length: %d[bytes]\n", d.Length())

	if _, err := io.Copy(a.player, d); err != nil {
		return err
	}
	return nil
}
