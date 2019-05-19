package republique

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/sirupsen/logrus"
)

type CompilerError struct {
	line     int
	filename string
	msg      string
}

func (e CompilerError) Error() string {
	return fmt.Sprintf("ERROR: Line %d: %s - %s", e.line, e.filename, e.msg)
}

type Compiler struct {
	log            *logrus.Logger
	filename       string
	outfile        string
	lines          []string
	indents        int
	command        *Command
	lastSubCommand *Command
}

func NewCompiler(log *logrus.Logger) *Compiler {
	return &Compiler{
		log:   log,
		lines: []string{},
	}
}

func (c *Compiler) Compile(filename string) error {
	c.log.WithField("filename", filename).Debug("Compiling")
	c.filename = filename
	ext := filepath.Ext(filename)
	c.outfile = filename[:len(filename)-len(ext)] + ".json"

	if err := c.load(); err != nil {
		c.log.WithError(err).WithField("filename", c.filename).Error("Failed to load file")
		return CompilerError{0, filename, "Loading: " + err.Error()}
	}

	switch ext {
	case ".oob":
		cmd, err := c.parseOOB()
		if err != nil {
			c.log.WithFields(logrus.Fields{
				"filename": c.filename,
				"numlines": len(c.lines),
			}).WithError(err).Debug("Failed to parse file")
			println(err.Error())
			return err
		}
		f, err := os.Create(c.outfile)
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
		// TODO
	case ".army":
		// TODO
	case ".game":
		// TODO
	}

	c.log.WithFields(logrus.Fields{
		"filename":   c.filename,
		"outputfile": c.outfile,
		"numlines":   len(c.lines),
		"indents":    c.indents,
	}).Debug("Loaded")

	return nil
}

func (c *Compiler) load() error {
	file, err := os.Open(c.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	c.lines = []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		c.lines = append(c.lines, scanner.Text())
	}
	return scanner.Err()
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
