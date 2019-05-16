package republique

var leaders = map[string]int64{
	// French Commanders
	"Augereau":        2,
	"Bernadotte":      1,
	"Bessieres":       2,
	"Brune":           1,
	"St Cyr":          2,
	"Clausel":         2,
	"Davout":          3,
	"Desaix":          3,
	"d'Erlon":         2,
	"Eugene":          2,
	"Friant":          2,
	"Foy":             2,
	"Gerard":          2,
	"Grouchy":         1,
	"Jerome":          1,
	"Jourdan":         2,
	"Junot":           2,
	"Kellerman":       2,
	"Kellerman Jr":    2,
	"Lannes":          3,
	"Lasalle":         2,
	"Latour Maubourg": 2,
	"Lefebvre":        2,
	"Legrand":         2,
	"MacDonald":       2,
	"Marmont":         2,
	"Massena":         3,
	"Moncey":          1,
	"Morand":          2,
	"Mortier":         2,
	"Murat":           4,
	"Napoleon":        5,
	"Ney":             4,
	"Oudinot":         3,
	"Poniatowski":     2,
	"Serurier":        2,
	"Soult":           3,
	"Suchet":          3,
	"Vandamme":        2,
	"Victor":          2,
	// British Commanders
	"Alten":            2,
	"Baird":            2,
	"Beresford":        2,
	"Colburn":          1,
	"Cole":             1,
	"Craufurd":         3,
	"Dalhousie":        1,
	"Fraser":           2,
	"Graham":           2,
	"Hill":             2,
	"Hope":             3,
	"Leith":            1,
	"Moore":            3,
	"Paget (Uxbridge)": 2,
	"Pakenham":         2,
	"Picton":           3,
	"Wellington":       4,
	// Prussian Commanders
	"Blucher":               4,
	"Brunswick":             2,
	"Friedrich-William III": 1,
	"Lestocq":               2,
	"Hohenlohe":             1,
	"Prince Louis":          2,
	"Zechwitz":              1,
	"Tauentzein":            1,
	"Prittwitz":             1,
	"von Schmettau":         1,
	"Wartensleben":          1,
	"Bulow":                 2,
	"Gneisnau":              2,
	"Scharnhorst":           3,
	"Thielmann":             2,
	"Yorck":                 3,
	"Ziethen":               2,
	// Russian Commanders
	"Alexander":        2,
	"Bagration":        3,
	"Barclay de Tolly": 2,
	"Bennigsen":        2,
	"Doctorov":         3,
	"Kutusov":          3,
	"Langeron":         2,
	"Miloradovitch":    2,
	"Platov":           3,
	"Raevski":          3,
	"Tolstoi":          1,
	"Tormasov":         2,
	"Tutchkov":         2,
	"Uvarov":           1,
	"Wittgenstein":     2,
	"Yermolov":         1,
}

func (c *Compiler) getLeaderRating() int64 {
	if c.command == nil {
		c.log.Error("No Commander, assigning 0 rating")
		return 0
	}
	m, ok := leaders[c.command.GetCommanderName()]
	if !ok {
		c.log.WithField("CommanderName", c.command.GetCommanderName()).Warn("Cannot find leader rating, assigning 1")
		return 1
	}
	return m
}
