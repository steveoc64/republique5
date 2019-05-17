package republique

import "strings"

var leaders = map[string]int64{
	// French Commanders
	"augereau":        2,
	"bernadotte":      1,
	"bessieres":       2,
	"brune":           1,
	"st cyr":          2,
	"clausel":         2,
	"davout":          3,
	"desaix":          3,
	"d'erlon":         2,
	"eugene":          2,
	"friant":          2,
	"foy":             2,
	"gerard":          2,
	"grouchy":         1,
	"jerome":          1,
	"jourdan":         2,
	"junot":           2,
	"kellerman":       2,
	"kellerman jr":    2,
	"lannes":          3,
	"lasalle":         2,
	"latour maubourg": 2,
	"lefebvre":        2,
	"legrand":         2,
	"macdonald":       2,
	"marmont":         2,
	"massena":         3,
	"moncey":          1,
	"morand":          2,
	"mortier":         2,
	"murat":           4,
	"napoleon":        5,
	"ney":             4,
	"oudinot":         3,
	"poniatowski":     2,
	"serurier":        2,
	"soult":           3,
	"suchet":          3,
	"vandamme":        2,
	"victor":          2,
	// British Commanders
	"alten":      2,
	"baird":      2,
	"beresford":  2,
	"colburn":    1,
	"cole":       1,
	"craufurd":   3,
	"dalhousie":  1,
	"fraser":     2,
	"graham":     2,
	"hill":       2,
	"hope":       3,
	"leith":      1,
	"moore":      3,
	"paget":      2,
	"uxbridge":   2,
	"pakenham":   2,
	"picton":     3,
	"wellington": 4,
	// Prussian Commanders
	"blucher":               4,
	"brunswick":             2,
	"friedrich-william iii": 1,
	"lestocq":               2,
	"hohenlohe":             1,
	"prince louis":          2,
	"zechwitz":              1,
	"tauentzein":            1,
	"prittwitz":             1,
	"von schmettau":         1,
	"wartensleben":          1,
	"bulow":                 2,
	"gneisnau":              2,
	"scharnhorst":           3,
	"thielmann":             2,
	"yorck":                 3,
	"ziethen":               2,
	// Russian Commanders
	"alexander":        2,
	"bagration":        3,
	"barclay de tolly": 2,
	"bennigsen":        2,
	"doctorov":         3,
	"kutusov":          3,
	"langeron":         2,
	"miloradovitch":    2,
	"platov":           3,
	"raevski":          3,
	"tolstoi":          1,
	"tormasov":         2,
	"tutchkov":         2,
	"uvarov":           1,
	"wittgenstein":     2,
	"yermolov":         1,
}

func (c *Compiler) getLeaderRating(name string) int64 {
	if c.command == nil {
		c.log.Warn("No Commander, assigning -1 rating")
		return -1
	}
	m, ok := leaders[strings.ToLower(name)]
	if !ok {
		c.log.WithField("CommanderName", name).Debug("Cannot find leader rating, assigning 0")
		return 0
	}
	return m
}
