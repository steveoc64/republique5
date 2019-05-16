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
	indents        int
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
		"filename":   c.filename,
		"outputfile": c.outfile,
		"numlines":   len(c.lines),
		"indents":    c.indents,
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
	year := 1800
	c.command = &Command{
		Arm:           Arm_INFANTRY,
		CommandRating: CommandRating_FUNCTIONAL,
		Nationality:   Nationality_ANY_NATION,
		Grade:         UnitGrade_REGULAR,
	}
	c.indents = 1
	var err error
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
				c.command.Arm = Arm_CAVALRY
			case "infantry":
				c.command.Arm = Arm_INFANTRY
			case "guards", "guard":
				c.command.Grade = UnitGrade_GUARD
			case "artillery":
				c.command.Arm = Arm_ARTILLERY
			case "elite":
				c.command.Grade = UnitGrade_ELITE
			case "veteran":
				c.command.Grade = UnitGrade_VETERAN
			case "regular":
				c.command.Grade = UnitGrade_REGULAR
			case "green", "conscript":
				c.command.Grade = UnitGrade_CONSCRIPT
			case "militia", "landwehr":
				c.command.Grade = UnitGrade_MILITIA
			case "rabble":
				c.command.Grade = UnitGrade_CIVILIAN
			case "efficient":
				c.command.CommandRating = CommandRating_EFFICIENT
			case "functional":
				c.command.CommandRating = CommandRating_FUNCTIONAL
			case "cumbersome":
				c.command.CommandRating = CommandRating_CUMBERSOME
			case "useless":
				c.command.CommandRating = CommandRating_USELESS
			case "french", "france":
				// TODO - turn this into a lambda function
				if ww != 2 {
					return k + 1, fmt.Errorf("!%s - missing year", w)
				}
				year, err = strconv.Atoi(words[1])
				if err != nil || year == 0 {
					return k + 1, fmt.Errorf("!%s - invalid year '%v'", w, words[1])
				}
				c.command.Nationality = Nationality_FRENCH
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
			cc := &Command{
				CommandRating: c.command.CommandRating,
				Arm:           c.command.Arm,
				Nationality:   c.command.Nationality,
				Grade:         c.command.Grade,
			}
			cc.Name = strings.TrimSpace(words[0])
			if ll == 2 {
				cc.CommanderName = strings.TrimSpace(words[1])
				cc.CommanderBonus = c.getLeaderRating(cc.CommanderName)
			}

			// Scan the title for rank strings
			cc.Rank = Rank_DIVISION
			lname := strings.ToLower(cc.Name)
			switch {
			case strings.Contains(lname, "cavalry div"),
				strings.Contains(lname, "cuirassier div"),
				strings.Contains(lname, "dragoon div"):
				cc.Arm = Arm_CAVALRY
				cc.Rank = Rank_CAVALRY_DIV
			case strings.Contains(lname, "cavalry brigade"):
				cc.Arm = Arm_CAVALRY
				cc.Rank = Rank_CAVALRY_BDE
			case strings.Contains(lname, "cavalry bde"),
				strings.Contains(lname, "hussar bde"),
				strings.Contains(lname, "chasseur bde"):
				cc.Arm = Arm_CAVALRY
				cc.Rank = Rank_CAVALRY_BDE
			case strings.Contains(lname, "artillery"):
				cc.Arm = Arm_ARTILLERY
				cc.Rank = Rank_GUN_PARK
			case strings.Contains(lname, "brigade"),
				strings.Contains(lname, "bde"):
				cc.Rank = Rank_BRIGADE
			}
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
		c.command.CommanderBonus = c.getLeaderRating(c.command.CommanderName)
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
