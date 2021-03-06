package compiler

import (
	"bufio"
	"fmt"
	"github.com/steveoc64/republique5/db"
	"os"
	"path/filepath"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/sirupsen/logrus"
)

// Error is an error type returned by the compiler suite
type Error struct {
	line     int
	filename string
	msg      string
}

// Error returns a string
func (e Error) Error() string {
	return fmt.Sprintf("ERROR: Line %d: %s - %s", e.line, e.filename, e.msg)
}

// Compiler is a holder struct for a compiler type
type Compiler struct {
	log *logrus.Logger
}

// NewCompiler returns a new compiler
func NewCompiler(log *logrus.Logger) *Compiler {
	return &Compiler{
		log: log,
	}
}

// Compile takes a file and calls the appropriate compiler depending on the filename extension
func (c *Compiler) Compile(filename string) error {
	ext := filepath.Ext(filename)
	shortName := filename[:len(filename)-len(ext)]

	switch ext {
	case ".oob":
		cmd, err := c.CompileOOB(filename)
		if err != nil {
			println(err.Error())
			return err
		}
		f, err := os.Create(shortName + ".json")
		if err != nil {
			fmt.Println(err)
			return err
		}
		marshaler := &jsonpb.Marshaler{}
		marshaler.Marshal(f, cmd)
		err = f.Close()
		if err != nil {
			fmt.Println(err)
			return err
		}
	case ".scenario":
		scn, err := c.CompileScenario(filename)
		if err != nil {
			println(err.Error())
			return err
		}
		f, err := os.Create(shortName + ".json")
		if err != nil {
			fmt.Println(err)
			return err
		}
		marshaler := &jsonpb.Marshaler{}
		marshaler.Marshal(f, scn)
		err = f.Close()
		if err != nil {
			fmt.Println(err)
			return err
		}
	case ".game":
		game, err := c.CompileGame(filename)
		if err != nil {
			println(err.Error())
			return err
		}

		db := db.NewDB(c.log, filepath.Base(shortName+".db"))
		err = db.Save("game", "state", game)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

func (c *Compiler) load(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func countLeadingRune(line string, r rune) int {
	i := 0
	for _, runeValue := range line {
		if runeValue == r {
			i++
		} else {
			break
		}
	}
	return i
}
