package compiler

import (
	"fmt"
	rp "github.com/steveoc64/republique5/proto"
	"runtime/debug"
	"strconv"
	"strings"
)

// CompileOOB reads an oob file and returns a compiled OOB
func (c *Compiler) CompileOOB(filename string) (*rp.Command, error) {
	lines, err := c.load(filename)
	if err != nil {
		return nil, fmt.Errorf("Error Loading %s: %s", filename, err.Error())
	}

	year := 1800
	skRating := rp.SkirmishRating_POOR
	skMax := "one"
	bnGuns := false
	cmd := &rp.Command{
		Arm:           rp.Arm_INFANTRY,
		CommandRating: rp.CommandRating_CUMBERSOME,
		Nationality:   rp.Nationality_ANY_NATION,
		Grade:         rp.UnitGrade_REGULAR,
		Drill:         rp.Drill_LINEAR,
	}
	var lastSubCommand *rp.Command
	indents := 1

	var k int
	var v string

	// catch panics
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ERROR:", r, "line", k+1, "file", filename, "\n->", v)
			debug.PrintStack()
		}
	}()

	getYear := func(k int, w []string) (int, error) {
		if k != 0 {
			return 0, CompilerError{k + 1, filename, "Nationality and Year must only be added on line 1"}
		}
		if len(w) != 2 {
			return 0, CompilerError{k + 1, filename, fmt.Sprintf("%s - missing year", strings.Join(w, " "))}
		}
		year, err = strconv.Atoi(w[1])
		if err != nil || year == 0 {
			return 0, CompilerError{k + 1, filename, fmt.Sprintf("%s - invalid year '%v'", strings.Join(w, " "), v)}
		}
		return year, nil
	}

	// scan for !commands
	for k, v = range lines {
		words := strings.Split(v, " ")
		ww := len(words)
		w := strings.ToLower(words[0])
		if len(v) == 0 {
			// empty lines
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
			case "cavalry":
				cmd.Arm = rp.Arm_CAVALRY
			case "infantry":
				cmd.Arm = rp.Arm_INFANTRY
			case "guards", "guard":
				cmd.Grade = rp.UnitGrade_GUARD
			case "artillery":
				cmd.Arm = rp.Arm_ARTILLERY
			case "elite":
				cmd.Grade = rp.UnitGrade_ELITE
			case "veteran":
				cmd.Grade = rp.UnitGrade_VETERAN
			case "regular":
				cmd.Grade = rp.UnitGrade_REGULAR
			case "green", "conscript":
				cmd.Grade = rp.UnitGrade_CONSCRIPT
			case "militia", "landwehr":
				cmd.Grade = rp.UnitGrade_MILITIA
			case "rabble":
				cmd.Grade = rp.UnitGrade_CIVILIAN
			case "efficient":
				cmd.CommandRating = rp.CommandRating_EFFICIENT
			case "functional":
				cmd.CommandRating = rp.CommandRating_FUNCTIONAL
			case "cumbersome":
				cmd.CommandRating = rp.CommandRating_CUMBERSOME
			case "useless":
				cmd.CommandRating = rp.CommandRating_USELESS
			case "linear":
				cmd.Drill = rp.Drill_LINEAR
			case "Massed":
				cmd.Drill = rp.Drill_MASSED
			case "Rapid":
				cmd.Drill = rp.Drill_RAPID
			case "french", "france":
				year, err = getYear(k, words)
				if err != nil {
					return nil, CompilerError{k + 1, filename, err.Error()}
				}
				cmd.Nationality = rp.Nationality_FRENCH
				switch {
				case year >= 1813:
					cmd.Drill = rp.Drill_RAPID
					cmd.CommandRating = rp.CommandRating_FUNCTIONAL
					cmd.Grade = rp.UnitGrade_CONSCRIPT
					skRating = rp.SkirmishRating_ADEQUATE
				case year >= 1805:
					cmd.Drill = rp.Drill_RAPID
					cmd.CommandRating = rp.CommandRating_EFFICIENT
					cmd.Grade = rp.UnitGrade_VETERAN
					skMax = "all"
					skRating = rp.SkirmishRating_CRACK_SHOT
				case year >= 1796:
					cmd.Drill = rp.Drill_MASSED
					cmd.CommandRating = rp.CommandRating_FUNCTIONAL
					cmd.Grade = rp.UnitGrade_VETERAN
					skRating = rp.SkirmishRating_CRACK_SHOT
				case year >= 1791:
					cmd.Drill = rp.Drill_MASSED
				}
			case "prussia", "prussian":
				year, err = getYear(k, words)
				if err != nil {
					return nil, CompilerError{k + 1, filename, fmt.Sprintf("Year '%v' %s", words, err.Error())}
				}
				cmd.Nationality = rp.Nationality_PRUSSIAN
				switch {
				case year >= 1814:
					cmd.Drill = rp.Drill_RAPID
					cmd.CommandRating = rp.CommandRating_EFFICIENT
					skRating = rp.SkirmishRating_ADEQUATE
				case year >= 1812:
					cmd.Drill = rp.Drill_MASSED
					cmd.CommandRating = rp.CommandRating_FUNCTIONAL
					cmd.Grade = rp.UnitGrade_CONSCRIPT
					skRating = rp.SkirmishRating_ADEQUATE
				case year <= 1806:
					bnGuns = true
				}
			case "austria", "austrian":
				year, err = getYear(k, words)
				if err != nil {
					return nil, CompilerError{k + 1, filename, fmt.Sprintf("Year '%v' %s", words, err.Error())}
				}
				cmd.Nationality = rp.Nationality_AUSTRIAN
				switch {
				case year >= 1813:
					cmd.Drill = rp.Drill_MASSED
					cmd.CommandRating = rp.CommandRating_FUNCTIONAL
					skRating = rp.SkirmishRating_ADEQUATE
				case year >= 1809:
					cmd.Drill = rp.Drill_MASSED
					skRating = rp.SkirmishRating_ADEQUATE
				case year <= 1802:
					bnGuns = true
				}
			case "russia", "russian":
				year, err = getYear(k, words)
				if err != nil {
					return nil, CompilerError{k + 1, filename, fmt.Sprintf("Year '%v' %s", words, err.Error())}
				}
				cmd.Nationality = rp.Nationality_RUSSIAN
				skMax = "none"
				switch {
				case year <= 1808:
					bnGuns = true
				}
			case "sweden":
				year, err = getYear(k, words)
				if err != nil {
					return nil, CompilerError{k + 1, filename, fmt.Sprintf("Year '%v' %s", words, err.Error())}
				}
				cmd.Nationality = rp.Nationality_SWEDEN
				skMax = "one"
				skRating = rp.SkirmishRating_ADEQUATE
				cmd.Drill = rp.Drill_LINEAR
				cmd.CommandRating = rp.CommandRating_FUNCTIONAL
			default:
				return nil, CompilerError{k + 1, filename, fmt.Sprintf("Invalid Command '%v'", v)}
			}
			// strip the line out
			//c.lines = append(c.lines[:k], c.lines[k+1:]...)
			continue
		}
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
		case 0: // Corps Definition
			words = strings.Split(v, " - ")
			if len(words) != 2 {
				return nil, CompilerError{k + 1, filename, "Invalid Corps Definition : needs 'Corps Name' - 'Commander Name'"}
			}
			params := words[1]
			ib1 := strings.Index(params, "(")
			ib2 := strings.Index(params, ")")
			if ib1 != -1 && ib2 != -1 {
				cmd.Notes = params[ib1+1 : ib2]
				params = params[:ib1]
			}
			cmd.Name = strings.TrimSpace(words[0])
			cmd.CommanderName = strings.TrimSpace(words[1])
			cmd.Rank = rp.Rank_CORPS
			cmd.Subcommands = []*rp.Command{}
			cmd.Units = []*rp.Unit{}
			cmd.CommanderBonus = -1
			if cmd.CommanderName != "" {
				cmd.CommanderBonus = c.getLeaderRating(cmd.CommanderName)
			}
			lastSubCommand = nil
			continue
		case 1: // Division SubCommand
			words = strings.Split(v[ioffset:], " - ")
			ll := len(words)
			if ll != 2 && ll != 1 {
				return nil, CompilerError{k + 1, filename, "Invalid Subcommand Definition - needs 'Subcommand Name' (- 'Commander Name')"}
			}
			subCommand := &rp.Command{
				CommandRating: cmd.CommandRating,
				Arm:           cmd.Arm,
				Nationality:   cmd.Nationality,
				Grade:         cmd.Grade,
				Drill:         cmd.Drill,
			}
			subCommand.Name = strings.TrimSpace(words[0])
			if strings.HasPrefix(strings.ToLower(subCommand.Name), "reserve") ||
				strings.HasPrefix(strings.ToLower(subCommand.Name), "bde reserve") ||
				strings.HasPrefix(strings.ToLower(subCommand.Name), "cavalry reserve") ||
				strings.HasPrefix(strings.ToLower(subCommand.Name), "div reserve") {
				subCommand.Reserve = true
			}
			params := ""
			if ll == 2 {
				params = words[1]
				ib1 := strings.Index(params, "(")
				ib2 := strings.Index(params, ")")
				if ib1 != -1 && ib2 != -1 {
					subCommand.Notes = params[ib1+1 : ib2]
					params = params[:ib1]
				}
				ib1 = strings.Index(params, "[")
				ib2 = strings.Index(params, "]")
				if ib1 != -1 && ib2 != -1 {
					nation := strings.ToLower(params[ib1+1 : ib2])
					switch nation {
					case "french":
						subCommand.Nationality = rp.Nationality_FRENCH
						skMax = "all"
					case "prussian", "saxon":
						subCommand.Nationality = rp.Nationality_PRUSSIAN
						skMax = "one"
					case "austrian":
						subCommand.Nationality = rp.Nationality_AUSTRIAN
						skMax = "one"
					case "russian":
						subCommand.Nationality = rp.Nationality_RUSSIAN
						skMax = "none"
					}
					params = params[:ib1]
				}
				ib1 = strings.Index(params, "[")

				subCommand.CommanderName = strings.TrimSpace(params)
				subCommand.CommanderBonus = c.getLeaderRating(subCommand.CommanderName)
			}

			// Scan the title for rank strings
			subCommand.Rank = rp.Rank_DIVISION
			lname := strings.ToLower(subCommand.Name)
			switch {
			case strings.Contains(lname, "cavalry div"),
				strings.Contains(lname, "cuirassier div"),
				strings.Contains(lname, "cossack div"),
				strings.Contains(lname, "corps cav"),
				strings.Contains(lname, "cavalry reserve"),
				strings.Contains(lname, "dragoon div"):
				subCommand.Arm = rp.Arm_CAVALRY
				subCommand.Rank = rp.Rank_CAVALRY_DIV
			case strings.Contains(lname, "cavalry brigade"):
				subCommand.Arm = rp.Arm_CAVALRY
				subCommand.Rank = rp.Rank_CAVALRY_BDE
			case strings.Contains(lname, "cavalry bde"),
				strings.Contains(lname, "hussar bde"),
				strings.Contains(lname, "chasseur bde"):
				subCommand.Arm = rp.Arm_CAVALRY
				subCommand.Rank = rp.Rank_CAVALRY_BDE
			case strings.Contains(lname, "artillery"):
				subCommand.Arm = rp.Arm_ARTILLERY
				subCommand.Rank = rp.Rank_GUN_PARK
			case strings.Contains(lname, "brigade"),
				strings.Contains(lname, "bde"):
				subCommand.Rank = rp.Rank_BRIGADE
			}
			subCommand.Units = []*rp.Unit{}
			lastSubCommand = subCommand
			cmd.Subcommands = append(cmd.Subcommands, subCommand)
			continue
		case 2: // Unit Definiition
			if lastSubCommand == nil {
				return nil, CompilerError{k + 1, filename, fmt.Sprintf("Indentation error - unit has no parent sub-command '%v'", v)}
			}
			hasGrade := false
			isGrenz := false
			v = strings.TrimSpace(v)
			words = strings.Split(v, " - ")
			if len(words) != 2 {
				return nil, CompilerError{k + 1, filename, "Invalid Unit Definition - needs 'Unit Name' - N bases [Unit Paramaters]"}
			}
			unit := &rp.Unit{
				Name:           strings.TrimSpace(words[0]),
				Arm:            lastSubCommand.Arm,
				UnitType:       rp.UnitType_INFANTRY_LINE,
				Grade:          lastSubCommand.Grade,
				Nationality:    lastSubCommand.Nationality,
				Drill:          lastSubCommand.Drill,
				CommandReserve: lastSubCommand.Reserve,
			}
			// max inherited grade is regular, except for guard formations which are all guard by default
			if unit.Grade > rp.UnitGrade_REGULAR && unit.Grade != rp.UnitGrade_GUARD {
				unit.Grade = rp.UnitGrade_REGULAR
			}

			params := words[1]
			pwords := strings.Split(params, " ")
			numBases, err := strconv.Atoi(pwords[0])
			if err != nil || numBases == 0 {
				numBases = 1
			} else {
				pwords = pwords[1:]
				// burn the next word if its bases
				switch pwords[0] {
				case "base", "bases":
					pwords = pwords[1:]
				}
			}
			unit.Strength = int32(numBases)
			// now join whats left back together
			params = strings.Join(pwords, " ")
			// extract notes if there are any
			ib1 := strings.Index(params, "(")
			ib2 := strings.Index(params, ")")
			if ib1 != -1 && ib2 != -1 {
				unit.Notes = params[ib1+1 : ib2]
				params = params[:ib1]
			}
			// now look for containing strings
			params = strings.ToLower(params)
			useSK := skRating
			useMax := skMax

			// gradings
			switch {
			case strings.Contains(params, "militia"),
				strings.Contains(params, "landwehr"):
				unit.Grade = rp.UnitGrade_MILITIA
				useSK = useSK.Decrement()
				useSK = useSK.Decrement()
				hasGrade = true
			case strings.Contains(params, "green"),
				strings.Contains(params, "conscript"):
				unit.Grade = rp.UnitGrade_CONSCRIPT
				useSK = useSK.Decrement()
				hasGrade = true
			case strings.Contains(params, "regular"):
				unit.Grade = rp.UnitGrade_REGULAR
				hasGrade = true
			case strings.Contains(params, "veteran"):
				unit.Grade = rp.UnitGrade_VETERAN
				hasGrade = true
			case strings.Contains(params, "elite"):
				unit.Grade = rp.UnitGrade_ELITE
				useSK = useSK.Increment()
				hasGrade = true
			case strings.Contains(params, "guard"):
				unit.Grade = rp.UnitGrade_GUARD
				useSK = useSK.Increment()
				useSK = useSK.Increment()
				hasGrade = true
			}

			// troop types
			switch {
			case strings.Contains(params, "rifle"):
				unit.Arm = rp.Arm_INFANTRY
				unit.UnitType = rp.UnitType_INFANTRY_RIFLES
				useSK = rp.SkirmishRating_EXCELLENT
			case strings.Contains(params, "grenadier"):
				unit.Arm = rp.Arm_INFANTRY
				unit.UnitType = rp.UnitType_INFANTRY_GRENADIER
			case strings.Contains(params, "dragoon"):
				unit.Arm = rp.Arm_CAVALRY
				unit.UnitType = rp.UnitType_CAVALRY_DRAGOON
				useMax = "all"
			case strings.Contains(params, "medium cav"):
				unit.Arm = rp.Arm_CAVALRY
				unit.UnitType = rp.UnitType_CAVALRY_MEDIUM
			case strings.Contains(params, "light cav"):
				unit.Arm = rp.Arm_CAVALRY
				unit.UnitType = rp.UnitType_CAVALRY_LIGHT
				useMax = "all"
			case strings.Contains(params, "heavy cav"):
				unit.Arm = rp.Arm_CAVALRY
				unit.UnitType = rp.UnitType_CAVALRY_HEAVY
			case strings.Contains(params, "hussar"):
				unit.Arm = rp.Arm_CAVALRY
				unit.UnitType = rp.UnitType_CAVALRY_HUSSAR
			case strings.Contains(params, "chas chev"),
				strings.Contains(params, "chaschev"),
				strings.Contains(params, "chev legere"),
				strings.Contains(params, "chasseur cheval"),
				strings.Contains(params, "chasseur a cheval"),
				strings.Contains(params, "chasseurs a'cheval"),
				strings.Contains(params, "chasseurs cheval"),
				strings.Contains(params, "chasseurs a cheval"),
				strings.Contains(params, "horse jager"),
				strings.Contains(params, "mounted jager"):
				unit.Arm = rp.Arm_CAVALRY
				unit.UnitType = rp.UnitType_CAVALRY_LIGHT
			case strings.Contains(params, "cuirassier"),
				strings.Contains(params, "carabinier"),
				strings.Contains(params, "karabinier"),
				strings.Contains(params, "kuirassier"):
				unit.Arm = rp.Arm_CAVALRY
				unit.UnitType = rp.UnitType_CAVALRY_CUIRASSIER
			case strings.Contains(params, "lancer"),
				strings.Contains(params, "landwehr cav"),
				strings.Contains(params, "uhlan"):
				unit.Arm = rp.Arm_CAVALRY
				unit.UnitType = rp.UnitType_CAVALRY_LANCER
			case strings.Contains(params, "cossack"):
				unit.Arm = rp.Arm_CAVALRY
				unit.UnitType = rp.UnitType_CAVALRY_COSSACK
			case strings.Contains(params, "mdf"):
				unit.Arm = rp.Arm_ARTILLERY
				unit.UnitType = rp.UnitType_ARTILLERY_MEDIUM
			case strings.Contains(params, "ltf"):
				unit.Arm = rp.Arm_ARTILLERY
				unit.UnitType = rp.UnitType_ARTILLERY_LIGHT
			case strings.Contains(params, "lth"):
				unit.Arm = rp.Arm_ARTILLERY
				unit.UnitType = rp.UnitType_ARTILLERY_LIGHT_HORSE
			case strings.Contains(params, "hvf"):
				unit.Arm = rp.Arm_ARTILLERY
				unit.UnitType = rp.UnitType_ARTILLERY_HEAVY
			case strings.Contains(params, "mdh"):
				unit.Arm = rp.Arm_ARTILLERY
				unit.UnitType = rp.UnitType_ARTILLERY_HORSE
			case strings.Contains(params, "light"),
				strings.Contains(params, "fusilier"),
				strings.Contains(params, "jager"):
				unit.Arm = rp.Arm_INFANTRY
				unit.UnitType = rp.UnitType_INFANTRY_LIGHT
				useMax = "all"
			case strings.Contains(params, "grenz"):
				unit.Arm = rp.Arm_INFANTRY
				unit.UnitType = rp.UnitType_INFANTRY_LIGHT
				useMax = "all"
				isGrenz = true
			case strings.Contains(params, "line"):
				unit.Arm = rp.Arm_INFANTRY
				unit.UnitType = rp.UnitType_INFANTRY_LINE
			}
			// set skirmisher ratings
			if unit.Arm == rp.Arm_INFANTRY ||
				unit.UnitType == rp.UnitType_CAVALRY_DRAGOON ||
				unit.UnitType == rp.UnitType_CAVALRY_LIGHT {
				if !hasGrade && unit.Grade < rp.UnitGrade_REGULAR {
					useSK = useSK.Decrement()
				}
				unit.SkirmishRating = useSK
				if unit.UnitType != rp.UnitType_INFANTRY_LIGHT &&
					unit.Grade < rp.UnitGrade_VETERAN &&
					useMax == "all" {
					useMax = "one"
				}
				switch useMax {
				case "one":
					unit.SkirmisherMax = 1
				case "all":
					unit.SkirmisherMax = unit.Strength
				case "none":
					unit.SkirmishRating = rp.SkirmishRating_POOR
				}
			}
			// Minimum default grading for cav and artillery
			if !hasGrade && (unit.Arm == rp.Arm_CAVALRY || unit.Arm == rp.Arm_ARTILLERY) {
				if unit.Grade < rp.UnitGrade_REGULAR {
					unit.Grade = rp.UnitGrade_REGULAR
				}
			}
			// default bnGuns for line troops by nationality
			if bnGuns &&
				(unit.UnitType == rp.UnitType_INFANTRY_LINE ||
					unit.UnitType == rp.UnitType_INFANTRY_GRENADIER) {
				unit.BnGuns = true
			}
			// grenzer rule
			if isGrenz && !hasGrade {
				unit.Grade = rp.UnitGrade_REGULAR
				unit.SkirmishRating = rp.SkirmishRating_POOR
			}
			lastSubCommand.Units = append(lastSubCommand.Units, unit)
			continue
		default:
			return nil, CompilerError{k + 1, filename, fmt.Sprintf("Dont know what to do with a unit at indent level %d", ii)}
		}
	}
	return cmd, nil
}
