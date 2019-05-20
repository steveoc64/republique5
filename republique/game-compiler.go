package republique

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
)

func NewAccessCode() string {
	a := make([]byte, 4, 4)
	for i := 0; i < 4; i++ {
		a[i] = byte('0' + rand.Intn(9))
	}
	return string(a)
}

func (c *Compiler) compileGame(filename string) (*Game, error) {
	lines, err := c.load(filename)
	if err != nil {
		return nil, fmt.Errorf("Error Loading %s: %s", filename, err.Error())
	}

	game := &Game{
		AccessCode: NewAccessCode(),
		TableX:     6,
		TableY:     4,
	}
	println("Compiling Game", filename, "AccessCode =", game.AccessCode)
	//var currentTeam *Team
	indents := 1

	var k int
	var v string

	// catch panics
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ERROR:", r, "line:", k+1, "file:", filename, "\n->", v)
			debug.PrintStack()
		}
	}()

	// scan for !commands
	for k, v = range lines {
		words := strings.Split(v, " ")
		ww := len(words)
		w := strings.ToLower(words[0])
		if len(v) == 0 {
			// empty lines are OK
			continue
		}
		if strings.HasPrefix(v, "#") {
			// is a comment
			continue
		}
		if strings.HasPrefix(v, "!") {
			w = w[1:]
			switch w {
			case "indent":
				if ww != 2 {
					return nil, CompilerError{k + 1, filename, "!Indent Command - missing size"}
				}
				i, err := strconv.Atoi(words[1])
				if err != nil || i < 1 {
					return nil, CompilerError{k + 1, filename, fmt.Sprintf("!Indent Command - invalid size '%v'", words[1])}
				}
				indents = i
				continue
			}
		}

		// work out the indent and base it off that
		ii := 0
		ioffset := 0
		if strings.HasPrefix(v, " ") {
			ioffset = countLeadingRune(v, ' ')
			ii = ioffset / indents
		}
		if strings.HasPrefix(v, "\t") {
			ioffset = countLeadingRune(v, '\t') / indents
			ii = ioffset
		}
		switch ii {
		case 0: // Directive
			switch w {
			case "scenario":
				if len(words) != 2 {
					return nil, CompilerError{k + 1, filename, "Invalid Scenario Name"}
				}
				filename := filepath.Join(filepath.Dir(filename), "..", "scenarios", words[1]+".scenario")
				scn, err := c.compileScenario(filename)
				if err != nil {
					return nil, err
				}
				game.Scenario = scn
			case "table":
				if len(words) != 2 {
					return nil, CompilerError{k + 1, filename, "expecting 'Table XxY' in feet"}
				}
				xpos := strings.Index(words[1], "x")
				if xpos < 1 {
					return nil, CompilerError{k + 1, filename, "expecting 'Table XxY' in feet"}
				}
				x, err := strconv.Atoi(words[1][:xpos])
				y, err2 := strconv.Atoi(words[1][xpos+1:])
				if err != nil || err2 != nil {
					return nil, CompilerError{k + 1, filename, "expecting 'Table X Y' in feet"}
				}
				game.TableX = int32(x)
				game.TableY = int32(y)
			default: // is a team name
				println("processing team", words[0])
			}
		case 1: // Player command
		default:
			return nil, CompilerError{k + 1, filename, fmt.Sprintf("Dont know what to do with a line at indent level %d '%v", ii, v)}
		}
	}
	return game, nil
}
