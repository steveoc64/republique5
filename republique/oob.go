package republique

import (
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"
)

func (c *Compiler) parseOOB() (*Command, error) {
	year := 1800
	skRating := SkirmishRating_POOR
	skMax := "one"
	bnGuns := false
	c.command = &Command{
		Arm:           Arm_INFANTRY,
		CommandRating: CommandRating_CUMBERSOME,
		Nationality:   Nationality_ANY_NATION,
		Grade:         UnitGrade_REGULAR,
		Drill:         Drill_LINEAR,
	}
	c.indents = 1
	var err error

	var k int
	var v string

	// catch panics
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ERROR:", r, "line", k, "file", c.filename, "\n->", v)
			debug.PrintStack()
		}
	}()

	getYear := func(k int, w []string) (int, error) {
		if k != 0 {
			return 0, CompilerError{k + 1, c.filename, "Nationality and Year must only be added on line 1"}
		}
		if len(w) != 2 {
			return 0, CompilerError{k + 1, c.filename, fmt.Sprintf("%s - missing year", strings.Join(w, " "))}
		}
		year, err = strconv.Atoi(w[1])
		if err != nil || year == 0 {
			return 0, CompilerError{k + 1, c.filename, fmt.Sprintf("%s - invalid year '%v'", strings.Join(w, " "))}
		}
		return year, nil
	}

	// scan for !commands
	for k, v = range c.lines {
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
					return nil, CompilerError{k + 1, c.filename, "!Indent Command - missing size"}
				}
				i, err := strconv.Atoi(words[1])
				if err != nil || i < 1 {
					return nil, CompilerError{k + 1, c.filename, fmt.Sprintf("!Indent Command - invalid size '%v'", words[1])}
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
			case "linear":
				c.command.Drill = Drill_LINEAR
			case "Massed":
				c.command.Drill = Drill_MASSED
			case "Rapid":
				c.command.Drill = Drill_RAPID
			case "french", "france":
				year, err = getYear(k, words)
				if err != nil {
					return nil, CompilerError{k + 1, c.filename, err.Error()}
				}
				c.command.Nationality = Nationality_FRENCH
				switch {
				case year >= 1813:
					c.command.Drill = Drill_RAPID
					c.command.CommandRating = CommandRating_FUNCTIONAL
					c.command.Grade = UnitGrade_CONSCRIPT
					skRating = SkirmishRating_ADEQUATE
				case year >= 1805:
					c.command.Drill = Drill_RAPID
					c.command.CommandRating = CommandRating_EFFICIENT
					c.command.Grade = UnitGrade_VETERAN
					skMax = "all"
					skRating = SkirmishRating_CRACK_SHOT
				case year >= 1796:
					c.command.Drill = Drill_MASSED
					c.command.CommandRating = CommandRating_FUNCTIONAL
					c.command.Grade = UnitGrade_VETERAN
					skRating = SkirmishRating_CRACK_SHOT
				case year >= 1791:
					c.command.Drill = Drill_MASSED
				}
			case "prussia", "prussian":
				year, err = getYear(k, words)
				if err != nil {
					return nil, CompilerError{k + 1, c.filename, fmt.Sprintf("Year '%v' %s", words, err.Error())}
				}
				c.command.Nationality = Nationality_PRUSSIAN
				switch {
				case year >= 1814:
					c.command.Drill = Drill_RAPID
					c.command.CommandRating = CommandRating_EFFICIENT
					skRating = SkirmishRating_ADEQUATE
				case year >= 1812:
					c.command.Drill = Drill_MASSED
					c.command.CommandRating = CommandRating_FUNCTIONAL
					c.command.Grade = UnitGrade_CONSCRIPT
					skRating = SkirmishRating_ADEQUATE
				case year <= 1806:
					bnGuns = true
				}
			case "austria", "austrian":
				year, err = getYear(k, words)
				if err != nil {
					return nil, CompilerError{k + 1, c.filename, fmt.Sprintf("Year '%v' %s", words, err.Error())}
				}
				c.command.Nationality = Nationality_AUSTRIAN
				switch {
				case year >= 1813:
					c.command.Drill = Drill_MASSED
					c.command.CommandRating = CommandRating_FUNCTIONAL
					skRating = SkirmishRating_ADEQUATE
				case year >= 1809:
					c.command.Drill = Drill_MASSED
					skRating = SkirmishRating_ADEQUATE
				case year <= 1802:
					bnGuns = true
				}
			case "russia", "russian":
				year, err = getYear(k, words)
				if err != nil {
					return nil, CompilerError{k + 1, c.filename, fmt.Sprintf("Year '%v' %s", words, err.Error())}
				}
				c.command.Nationality = Nationality_RUSSIAN
				skMax = "none"
				switch {
				case year <= 1808:
					bnGuns = true
				}
			case "sweden":
				year, err = getYear(k, words)
				if err != nil {
					return nil, CompilerError{k + 1, c.filename, fmt.Sprintf("Year '%v' %s", words, err.Error())}
				}
				c.command.Nationality = Nationality_SWEDEN
				skMax = "one"
				skRating = SkirmishRating_ADEQUATE
				c.command.Drill = Drill_LINEAR
				c.command.CommandRating = CommandRating_FUNCTIONAL
			default:
				return nil, CompilerError{k + 1, c.filename, fmt.Sprintf("Invalid Command '%v'", v)}
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
		case 0: // Corps Definition
			words = strings.Split(v, " - ")
			if len(words) != 2 {
				return nil, CompilerError{k + 1, c.filename, "Invalid Corps Definition : needs 'Corps Name' - 'Commander Name'"}
			}
			params := words[1]
			ib1 := strings.Index(params, "(")
			ib2 := strings.Index(params, ")")
			if ib1 != -1 && ib2 != -1 {
				c.command.Notes = params[ib1+1 : ib2]
				params = params[:ib1]
			}
			c.command.Name = strings.TrimSpace(words[0])
			c.command.CommanderName = strings.TrimSpace(words[1])
			c.command.Rank = Rank_CORPS
			c.command.Subcommands = []*Command{}
			c.command.Units = []*Unit{}
			c.command.CommanderBonus = c.getLeaderRating(c.command.CommanderName)
			c.lastSubCommand = nil
			continue
		case 1: // Division SubCommand
			words = strings.Split(v[ioffset:], " - ")
			ll := len(words)
			if ll != 2 && ll != 1 {
				return nil, CompilerError{k + 1, c.filename, "Invalid Subcommand Definition - needs 'Subcommand Name' (- 'Commander Name')"}
			}
			cc := &Command{
				CommandRating: c.command.CommandRating,
				Arm:           c.command.Arm,
				Nationality:   c.command.Nationality,
				Grade:         c.command.Grade,
				Drill:         c.command.Drill,
			}
			cc.Name = strings.TrimSpace(words[0])
			if strings.HasPrefix(strings.ToLower(cc.Name), "reserve") ||
				strings.HasPrefix(strings.ToLower(cc.Name), "bde reserve") ||
				strings.HasPrefix(strings.ToLower(cc.Name), "cavalry reserve") ||
				strings.HasPrefix(strings.ToLower(cc.Name), "div reserve") {
				cc.Reserve = true
			}
			params := ""
			if ll == 2 {
				params = words[1]
				ib1 := strings.Index(params, "(")
				ib2 := strings.Index(params, ")")
				if ib1 != -1 && ib2 != -1 {
					cc.Notes = params[ib1+1 : ib2]
					params = params[:ib1]
				}
				ib1 = strings.Index(params, "[")
				ib2 = strings.Index(params, "]")
				if ib1 != -1 && ib2 != -1 {
					nation := strings.ToLower(params[ib1+1 : ib2])
					switch nation {
					case "french":
						cc.Nationality = Nationality_FRENCH
						skMax = "all"
					case "prussian", "saxon":
						cc.Nationality = Nationality_PRUSSIAN
						skMax = "one"
					case "austrian":
						cc.Nationality = Nationality_AUSTRIAN
						skMax = "one"
					case "russian":
						cc.Nationality = Nationality_RUSSIAN
						skMax = "none"
					}
					params = params[:ib1]
				}
				ib1 = strings.Index(params, "[")

				cc.CommanderName = strings.TrimSpace(params)
				cc.CommanderBonus = c.getLeaderRating(cc.CommanderName)
			}

			// Scan the title for rank strings
			cc.Rank = Rank_DIVISION
			lname := strings.ToLower(cc.Name)
			switch {
			case strings.Contains(lname, "cavalry div"),
				strings.Contains(lname, "cuirassier div"),
				strings.Contains(lname, "cossack div"),
				strings.Contains(lname, "corps cav"),
				strings.Contains(lname, "cavalry reserve"),
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
			continue
		case 2: // Unit Definiition
			if c.lastSubCommand == nil {
				return nil, CompilerError{k + 1, c.filename, fmt.Sprintf("Indentation error - unit has no parent sub-command '%v'", v)}
			}
			hasGrade := false
			isGrenz := false
			v = strings.TrimSpace(v)
			words = strings.Split(v, " - ")
			if len(words) != 2 {
				return nil, CompilerError{k + 1, c.filename, "Invalid Unit Definition - needs 'Unit Name' - N bases [Unit Paramaters]"}
			}
			unit := &Unit{
				Name:           strings.TrimSpace(words[0]),
				Arm:            c.lastSubCommand.Arm,
				UnitType:       UnitType_INFANTRY_LINE,
				Grade:          c.lastSubCommand.Grade,
				Nationality:    c.lastSubCommand.Nationality,
				Drill:          c.lastSubCommand.Drill,
				CommandReserve: c.lastSubCommand.Reserve,
			}
			// max inherited grade is regular, except for guard formations which are all guard by default
			if unit.Grade > UnitGrade_REGULAR && unit.Grade != UnitGrade_GUARD {
				unit.Grade = UnitGrade_REGULAR
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
			unit.Strength = int64(numBases)
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
				unit.Grade = UnitGrade_MILITIA
				useSK = useSK.Decrement()
				useSK = useSK.Decrement()
				hasGrade = true
			case strings.Contains(params, "green"),
				strings.Contains(params, "conscript"):
				unit.Grade = UnitGrade_CONSCRIPT
				useSK = useSK.Decrement()
				hasGrade = true
			case strings.Contains(params, "regular"):
				unit.Grade = UnitGrade_REGULAR
				hasGrade = true
			case strings.Contains(params, "veteran"):
				unit.Grade = UnitGrade_VETERAN
				hasGrade = true
			case strings.Contains(params, "elite"):
				unit.Grade = UnitGrade_ELITE
				useSK = useSK.Increment()
				hasGrade = true
			case strings.Contains(params, "guard"):
				unit.Grade = UnitGrade_GUARD
				useSK = useSK.Increment()
				useSK = useSK.Increment()
				hasGrade = true
			}

			// troop types
			switch {
			case strings.Contains(params, "rifle"):
				unit.Arm = Arm_INFANTRY
				unit.UnitType = UnitType_INFANTRY_RIFLES
				useSK = SkirmishRating_EXCELLENT
			case strings.Contains(params, "grenadier"):
				unit.Arm = Arm_INFANTRY
				unit.UnitType = UnitType_INFANTRY_GRENADIER
			case strings.Contains(params, "dragoon"):
				unit.Arm = Arm_CAVALRY
				unit.UnitType = UnitType_CAVALRY_DRAGOON
				useMax = "all"
			case strings.Contains(params, "medium cav"):
				unit.Arm = Arm_CAVALRY
				unit.UnitType = UnitType_CAVALRY_MEDIUM
			case strings.Contains(params, "light cav"):
				unit.Arm = Arm_CAVALRY
				unit.UnitType = UnitType_CAVALRY_LIGHT
				useMax = "all"
			case strings.Contains(params, "heavy cav"):
				unit.Arm = Arm_CAVALRY
				unit.UnitType = UnitType_CAVALRY_HEAVY
			case strings.Contains(params, "hussar"):
				unit.Arm = Arm_CAVALRY
				unit.UnitType = UnitType_CAVALRY_HUSSAR
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
				unit.Arm = Arm_CAVALRY
				unit.UnitType = UnitType_CAVALRY_LIGHT
			case strings.Contains(params, "cuirassier"),
				strings.Contains(params, "carabinier"),
				strings.Contains(params, "karabinier"),
				strings.Contains(params, "kuirassier"):
				unit.Arm = Arm_CAVALRY
				unit.UnitType = UnitType_CAVALRY_CUIRASSIER
			case strings.Contains(params, "lancer"),
				strings.Contains(params, "landwehr cav"),
				strings.Contains(params, "uhlan"):
				unit.Arm = Arm_CAVALRY
				unit.UnitType = UnitType_CAVALRY_LANCER
			case strings.Contains(params, "cossack"):
				unit.Arm = Arm_CAVALRY
				unit.UnitType = UnitType_CAVALRY_COSSACK
			case strings.Contains(params, "mdf"):
				unit.Arm = Arm_ARTILLERY
				unit.UnitType = UnitType_ARTILLERY_MEDIUM
			case strings.Contains(params, "ltf"):
				unit.Arm = Arm_ARTILLERY
				unit.UnitType = UnitType_ARTILLERY_LIGHT
			case strings.Contains(params, "lth"):
				unit.Arm = Arm_ARTILLERY
				unit.UnitType = UnitType_ARTILLERY_LIGHT_HORSE
			case strings.Contains(params, "hvf"):
				unit.Arm = Arm_ARTILLERY
				unit.UnitType = UnitType_ARTILLERY_HEAVY
			case strings.Contains(params, "mdh"):
				unit.Arm = Arm_ARTILLERY
				unit.UnitType = UnitType_ARTILLERY_HORSE
			case strings.Contains(params, "light"),
				strings.Contains(params, "fusiliers"),
				strings.Contains(params, "jager"):
				unit.Arm = Arm_INFANTRY
				unit.UnitType = UnitType_INFANTRY_LIGHT
				useMax = "all"
			case strings.Contains(params, "grenz"):
				unit.Arm = Arm_INFANTRY
				unit.UnitType = UnitType_INFANTRY_LIGHT
				useMax = "all"
				isGrenz = true
			case strings.Contains(params, "line"):
				unit.Arm = Arm_INFANTRY
				unit.UnitType = UnitType_INFANTRY_LINE
			}
			// set skirmisher ratings
			if unit.Arm == Arm_INFANTRY ||
				unit.UnitType == UnitType_CAVALRY_DRAGOON ||
				unit.UnitType == UnitType_CAVALRY_LIGHT {
				if !hasGrade && unit.Grade < UnitGrade_REGULAR {
					useSK = useSK.Decrement()
				}
				unit.SkirmishRating = useSK
				if unit.UnitType != UnitType_INFANTRY_LIGHT &&
					unit.Grade < UnitGrade_VETERAN &&
					useMax == "all" {
					useMax = "one"
				}
				switch useMax {
				case "one":
					unit.SkirmisherMax = 1
				case "all":
					unit.SkirmisherMax = unit.Strength
				case "none":
					unit.SkirmishRating = SkirmishRating_POOR
				}
			}
			// Minimum default grading for cav and artillery
			if !hasGrade && (unit.Arm == Arm_CAVALRY || unit.Arm == Arm_ARTILLERY) {
				if unit.Grade < UnitGrade_REGULAR {
					unit.Grade = UnitGrade_REGULAR
				}
			}
			// default bnGuns for line troops by nationality
			if bnGuns &&
				(unit.UnitType == UnitType_INFANTRY_LINE ||
					unit.UnitType == UnitType_INFANTRY_GRENADIER) {
				unit.BnGuns = true
			}
			// grenzer rule
			if isGrenz && !hasGrade {
				unit.Grade = UnitGrade_REGULAR
				unit.SkirmishRating = SkirmishRating_POOR
			}
			c.lastSubCommand.Units = append(c.lastSubCommand.Units, unit)
			continue
		default:
			return nil, CompilerError{k + 1, c.filename, fmt.Sprintf("Dont know what to do with a unit at indent level %d", ii)}
		}
	}
	return c.command, nil
}
