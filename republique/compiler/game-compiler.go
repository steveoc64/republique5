package compiler

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	rp "github.com/steveoc64/republique5/proto"
)

// NewAccessCode generates a random access code and returns it as a string
func NewAccessCode() string {
	a := make([]byte, 4, 4)
	for i := 0; i < 4; i++ {
		a[i] = byte('0' + rand.Intn(9))
	}
	return string(a)
}

// CompileGame reads a game file and returns a compiled game
func (c *Compiler) CompileGame(filename string) (*rp.Game, error) {
	rand.Seed(time.Now().UnixNano())
	lines, err := c.load(filename)
	if err != nil {
		return nil, fmt.Errorf("Error Loading %s: %s", filename, err.Error())
	}

	game := &rp.Game{
		AdminAccess: NewAccessCode(),
		TableX:      6,
		TableY:      4,
	}
	tableMode := false
	tableLine := int32(0)
	println("Compiling Game", filename, "AdminAccess =", game.AdminAccess)
	var currentTeam *rp.Team
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
		if tableMode {
			if tableLine == 0 && v[0] != '-' {
				tableMode = false
			} else {
				if tableLine > 0 && tableLine <= game.TableY {
					// v contains a row to add to the tableLayout
					for int32(len(v)) < game.TableX {
						v = v + " "
					}
					game.TableLayout = game.TableLayout + v
				}
				tableLine++
				if tableLine > game.TableY+1 {
					tableMode = false
				}
				continue
			}
		}
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
			case "name":
				game.Name = strings.TrimSpace(v[4:])
			case "scenario":
				if len(words) != 2 {
					return nil, CompilerError{k + 1, filename, "Invalid Scenario Name"}
				}
				filename := filepath.Join(filepath.Dir(filename), "..", "scenarios", words[1]+".scenario")
				scn, err := c.CompileScenario(filename)
				if err != nil {
					return nil, err
				}
				game.Scenario = scn
				game.GameTime = scn.StartTime
				for _, v := range scn.Teams {
					v.AccessCode = NewAccessCode()
					println("Team", v.Name, "AccessCode =", v.AccessCode)
				}
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
				tableMode = true
				// following this is a map
				// line of dashes + Y lines of X chars + line of dashes
			default: // is a team name
				currentTeam = nil
				teamname := words[0]
				teamwords := strings.Split(v, " - ")
				teamgamename := game.Name
				teamside := rp.MapSide_FRONT
				if len(teamwords) != 2 {
					return nil, CompilerError{k + 1,
						filename,
						fmt.Sprintf("invalid team name '%v' : expecting 'Team Name - Game Descriptnion (MapSide)' one of (Front, Top, Left, Right)", v)}
				}
				teamname = teamwords[0]
				teamgamename = teamwords[1]
				// extract the mapside
				l1 := strings.Index(teamgamename, "(")
				l2 := strings.Index(teamgamename, ")")
				if l1 == -1 || l2 == -1 {
					return nil, CompilerError{k + 1,
						filename,
						fmt.Sprintf("invalid team name '%v' : expecting 'Team Name - Game Descriptnion (MapSide)' one of (Front, Top, Left, Right)", v)}
				}
				sidename := strings.ToLower(teamgamename[l1+1 : l2])
				switch sidename {
				case "front":
					teamside = rp.MapSide_FRONT
				case "top":
					teamside = rp.MapSide_TOP
				case "left":
					teamside = rp.MapSide_LEFT_FLANK
				case "right":
					teamside = rp.MapSide_RIGHT_FLANK
				default:
					return nil, CompilerError{k + 1,
						filename,
						fmt.Sprintf("invalid team name '%v' : expecting 'Team Name - Game Descriptnion (MapSide)' one of (Front, Top, Left, Right)", v)}
				}
				teamgamename = teamgamename[:l1-1]
				for _, team := range game.Scenario.Teams {
					if team.Name == teamname {
						currentTeam = team
						team.GameName = teamgamename
						team.Side = teamside
						break
					}
				}
				if currentTeam == nil {
					return nil, CompilerError{k + 1, filename, fmt.Sprintf("Unknown Team '%v'", words[0])}
				}
				println("processing team", currentTeam.Name)
			}
		case 1: // Player command
			if currentTeam == nil {
				return nil, CompilerError{k + 1, filename, "No Team Defined"}
			}
			player := &rp.Player{AccessCode: NewAccessCode()}
			v = strings.TrimSpace(v)
			names := strings.Split(v, ",")
			for _, name := range names {
				name = strings.TrimSpace(name)
				cc := currentTeam.GetCommandByName(name)
				if cc == nil {
					cc = currentTeam.GetCommandByCommanderName(name)
					if cc == nil {
						return nil, CompilerError{k + 1, filename, fmt.Sprintf("%v is not a valid unit or commander in team %v", name, currentTeam.Name)}
					}
				}
				player.Commanders = append(player.Commanders, cc.CommanderName)
			}
			currentTeam.Players = append(currentTeam.Players, player)
		default:
			return nil, CompilerError{k + 1, filename, fmt.Sprintf("Dont know what to do with a line at indent level %d '%v", ii, v)}
		}
	}

	game.GenerateIDs()
	return game, nil
}
