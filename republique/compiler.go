package republique

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/gogo/protobuf/jsonpb"

	"github.com/sirupsen/logrus"
)

type Compiler struct {
	log            *logrus.Logger
	filename       string
	outfile        string
	lines          []string
	basedFrom      string
	indents        int
	grade          string
	arm            string
	cmdDoctrine    string
	command        *Command
	mode           int
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
		return err
	}

	switch ext {
	case ".oob":
		if i, err := c.parseOOB(); err != nil {
			c.log.WithFields(logrus.Fields{
				"filename":   c.filename,
				"numlines":   len(c.lines),
				"LineNumber": i,
			}).WithError(err).Error("Failed to parse file")
			return err
		}
	case "scenario":
		// TODO
	case "army":
		// TODO
	case "game":
		// TODO
	}

	c.log.WithFields(logrus.Fields{
		"filename":    c.filename,
		"outputfile":  c.outfile,
		"numlines":    len(c.lines),
		"indents":     c.indents,
		"basedFrom":   c.basedFrom,
		"cmdDoctrine": c.cmdDoctrine,
		"grade":       c.grade,
		"arm":         c.arm,
	}).Debug("Loaded")

	j := &bytes.Buffer{}
	marshaler := &jsonpb.Marshaler{}
	marshaler.Marshal(j, c.command)
	ioutil.WriteFile(c.outfile, j.Bytes(), 0644)

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

func (c *Compiler) parseOOB() (int, error) {
	c.grade = "Regular"
	c.arm = "Infantry"
	c.cmdDoctrine = "Functional"
	c.command = &Command{}
	c.indents = 1
	// scan for !commands
	for k, v := range c.lines {
		words := strings.Split(v, " ")
		ww := len(words)
		w := strings.ToLower(words[0])
		if len(v) == 0 {
			//c.lines = append(c.lines[:k], c.lines[k+1:]...)
			continue
		}
		if strings.HasPrefix(v, "#") {
			// is a comment
			//c.lines = append(c.lines[:k], c.lines[k+1:]...)
			continue
		}
		if strings.HasPrefix(v, "!") {
			w = w[1:]
			switch w {
			case "indent":
				if ww != 2 {
					return k + 1, fmt.Errorf("!Indent Command - missing size")
				}
				i, err := strconv.Atoi(words[1])
				if err != nil || i < 1 {
					return k + 1, fmt.Errorf("!Indent Command - invalid size '%v'", words[1])
				}
				c.indents = i
			case "cavalry":
				c.arm = "Cavalry"
			case "infantry":
				c.arm = "Infantry"
			case "guards", "guard":
				c.grade = "Guard"
			case "artillery":
				c.arm = "Artillery"
			case "elite":
				c.grade = "Elite"
			case "veteran":
				c.grade = "Veteran"
			case "regular":
				c.grade = "Regular"
			case "green", "conscript":
				c.grade = "Conscript"
			case "militia", "landwehr":
				c.grade = "Militia"
			case "rabble":
				c.grade = "Rabble"
			case "efficient":
				c.cmdDoctrine = "Efficient"
			case "functional":
				c.cmdDoctrine = "Functional"
			case "cumbersome":
				c.cmdDoctrine = "Cumbersome"
			case "french", "france":
				// TODO - turn this into a lambda function
				if ww != 2 {
					return k + 1, fmt.Errorf("!%s - missing year", w)
				}
				year, err := strconv.Atoi(words[1])
				if err != nil || year == 0 {
					return k + 1, fmt.Errorf("!%s - invalid year '%v'", w, words[1])
				}
				c.basedFrom = fmt.Sprintf("French-%d", year)
			case "austrian", "austria":
				c.basedFrom = fmt.Sprintf("Austrian-%d", 0)
			case "russian", "russia":
				c.basedFrom = fmt.Sprintf("Russian-%d", 0)
			case "prussian", "prussia":
				c.basedFrom = fmt.Sprintf("Prussian-%d", 0)
			case "british", "britain":
				c.basedFrom = fmt.Sprintf("British-%d", 0)
			case "spanish", "spain":
				c.basedFrom = fmt.Sprintf("Spanish-%d", 0)
			case "portuguese", "portugal":
				c.basedFrom = fmt.Sprintf("Spanish-%d", 0)
			case "ottoman", "turkish":
				c.basedFrom = fmt.Sprintf("Ottoman-%d", 0)
			case "dutch", "belgian":
				c.basedFrom = fmt.Sprintf("DutchBelgian-%d", 0)
			case "italian", "italy":
				c.basedFrom = fmt.Sprintf("Italian-%d", 0)
			case "persian", "persia":
				c.basedFrom = fmt.Sprintf("Persian-%d", 0)
			case "bavarian", "bavaria":
				c.basedFrom = fmt.Sprintf("Bavarian-%d", 0)
			case "american", "america":
				c.basedFrom = fmt.Sprintf("American-%d", 0)
			case "irregular":
				c.basedFrom = fmt.Sprintf("Irregular-%d", 0)
			case "native":
				c.basedFrom = fmt.Sprintf("Native-%d", 0)
			case "german":
				c.basedFrom = fmt.Sprintf("German-%d", 0)
			case "nassau":
				c.basedFrom = fmt.Sprintf("Nassau-%d", 0)
			default:
				return k + 1, fmt.Errorf("Invalid Command '%s'", v)
			}
			// strip the line out
			//c.lines = append(c.lines[:k], c.lines[k+1:]...)
			continue
		}
		ii := 0
		ioffset := 0
		if strings.HasPrefix(v, " ") {
			ioffset = countLeadingRune(v, ' ')
			ii = ioffset / c.indents
		}
		if strings.HasPrefix(v, "\t") {
			ioffset = countLeadingRune(v, '\t') / c.indents
			ii = ioffset
		}
		switch ii {
		case 0:
			//println(k+1, "looks like a corps definition line", v)
			c.mode = 1
		case 1:
			//println(k+1, "SubCommand:", v)
			words = strings.Split(v[ioffset:], "-")
			ll := len(words)
			if ll != 2 && ll != 1 {
				spew.Dump(ll, words)
				return k + 1, fmt.Errorf("Invalid Subcommand Definition - needs 'Subcommand Name' (- 'Commander Name')")
			}
			cc := &Command{}
			cc.Name = strings.TrimSpace(words[0])
			if ll == 2 {
				cc.CommanderName = strings.TrimSpace(words[1])
				cc.CommandRating = c.getLeaderRating()
			}
			cc.Rank = Rank_DIVISION
			cc.Subcommands = []*Command{}
			cc.Units = []*Unit{}
			c.lastSubCommand = cc
			c.command.Subcommands = append(c.command.Subcommands, cc)
			c.mode = 2
			continue
		case 2:
			//println(k+1, "Unit:", v)
			c.mode = 2
			continue
		default:
			return k + 1, fmt.Errorf("Dont know what to do with a unit at indent level %d", ii)
		}
		words = strings.Split(v, "-")
		if len(words) != 2 {
			return k + 1, fmt.Errorf("Invalid Corps Definition - needs 'Corps Name' - 'Commander Name'")
		}
		c.command.Name = strings.TrimSpace(words[0])
		c.command.CommanderName = strings.TrimSpace(words[1])
		c.command.Rank = Rank_CORPS
		c.command.Subcommands = []*Command{}
		c.command.Units = []*Unit{}
		c.command.CommandRating = c.getLeaderRating()
		c.lastSubCommand = nil
	}
	return 0, nil
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
