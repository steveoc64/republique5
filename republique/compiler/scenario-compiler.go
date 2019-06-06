package compiler

import (
	"fmt"
	rp "github.com/steveoc64/republique5/proto"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	tspb "github.com/golang/protobuf/ptypes/timestamp"
)

// ComplieScenario reads a scenario file and returns a compiled scenario
func (c *Compiler) CompileScenario(filename string) (*rp.Scenario, error) {
	lines, err := c.load(filename)
	if err != nil {
		return nil, fmt.Errorf("Error Loading %s: %s", filename, err.Error())
	}

	scn := &rp.Scenario{
		Teams: map[string]*rp.Team{},
	}
	var currentTeam *rp.Team
	indents := 1
	isBriefing := false
	position := rp.BattlefieldPosition_REAR

	var k int
	var v string

	// catch panics
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ERROR:", r, "line:", k+1, "file:", filename, "\n->", v)
			debug.PrintStack()
		}
	}()

	getFromTo := func(arr string) (int, int) {
		w := strings.Split(arr, " ")
		ww := strings.Split(w[0], "-")
		f, _ := strconv.Atoi(ww[0])
		t, _ := strconv.Atoi(ww[1])
		return f, t
	}
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
			case "start":
				timestring := v[7:]
				t, err := time.Parse("02 Jan 2006 15:04", strings.TrimSpace(timestring))
				if err != nil {
					return nil, CompilerError{k + 1, filename, fmt.Sprintf("Error parsing datetime '%v': %s", words[1], err.Error())}
				}
				scn.StartTime = &tspb.Timestamp{Seconds: t.Unix()}
				continue
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
		case 0: // Side Definition
			currentTeam = &rp.Team{
				Name: strings.TrimSpace(v),
			}
			scn.Teams[currentTeam.Name] = currentTeam
		case 1: // Player command
			if currentTeam == nil {
				return nil, CompilerError{k + 1, filename, "No player side defined at 1 indent above this line"}
			}
			command := strings.TrimSpace(v)
			switch strings.ToLower(command) {
			case "briefing":
				isBriefing = true
			case "center", "centre":
				isBriefing = false
				position = rp.BattlefieldPosition_CENTRE
			case "right":
				isBriefing = false
				position = rp.BattlefieldPosition_RIGHT
			case "left":
				isBriefing = false
				position = rp.BattlefieldPosition_LEFT
			case "reserve":
				isBriefing = false
				position = rp.BattlefieldPosition_REAR
			default:
				return nil, CompilerError{k + 1, filename, "Invalid command: " + v}
			}
			continue
		case 2: // unit or briefing
			if currentTeam == nil {
				return nil, CompilerError{k + 1, filename, "No player side defined at 1 indent above this line"}
			}
			if isBriefing {
				if currentTeam.Briefing != "" {
					currentTeam.Briefing = currentTeam.Briefing + "\n"
				}
				currentTeam.Briefing = currentTeam.Briefing + strings.TrimSpace(v)
				continue
			}
			// if we are still here, its a unit to be added to current player
			// extract out the arrival data
			v := strings.TrimSpace(v)
			ib1 := strings.Index(v, " (")
			ib2 := strings.Index(v, ")")
			arrival := ""
			chance := 100
			if ib1 != -1 && ib2 != -1 {
				arrival = strings.ToLower(v[ib1+2 : ib2])
				pc := strings.TrimSpace(v[ib2+1:])
				if strings.HasSuffix(pc, "%") {
					chance, _ = strconv.Atoi(pc[:strings.Index(pc, "%")])
				}
				v = v[:ib1]
			}
			fromHour := 0
			toHour := 0
			contact := false
			switch {
			case arrival == "contact": // start of game, but in contact
				contact = true
			case arrival == "": // start of game, but in reserve
			case strings.HasSuffix(arrival, " hrs"),
				strings.HasSuffix(arrival, " hours"):
				fromHour, toHour = getFromTo(arrival)
				//case strings.HasPrefix(arrival, "at "):
				//timestring := arrival[3:]
				//t, _ := time.Parse("02 Jan 2006 15:04", strings.TrimSpace(timestring))
				//tt := &tspb.Timestamp{Seconds: t.Unix()}
			}
			// Get the command unit
			w := strings.Split(v, " - ")
			subUnit := ""
			commandName := ""
			switch len(w) {
			case 1:
				commandName = strings.TrimSpace(v)
			case 2:
				commandName = strings.TrimSpace(w[0])
				subUnit = strings.TrimSpace(w[1])
			default:
				return nil, CompilerError{k + 1, filename, "Too many - chars, expecting 'MainOOB - SubUnit (arrivals) XX%"}
			}
			cmd, err := c.CompileOOB(filepath.Join(filepath.Dir(filename), commandName+".oob"))
			if err != nil {
				return nil, err
			}
			if subUnit != "" {
				found := false
				for _, v := range cmd.Subcommands {
					if v.Name == subUnit || v.CommanderName == subUnit {
						cmd = v
						found = true
						break
					}
				}
				if !found {
					return nil, CompilerError{k + 1, filename, fmt.Sprintf("Failed to find subUnit '%v' in '%v", subUnit, commandName)}
				}
			}
			cmd.Arrival = &rp.Arrival{
				From:     int32(fromHour),
				To:       int32(toHour),
				Percent:  int32(chance),
				Position: position,
				Contact:  contact,
			}
			currentTeam.Commands = append(currentTeam.Commands, cmd)

			continue
		default:
			return nil, CompilerError{k + 1, filename, fmt.Sprintf("Dont know what to do with a unit at indent level %d '%v", ii, v)}
		}
	}
	return scn, nil
}
