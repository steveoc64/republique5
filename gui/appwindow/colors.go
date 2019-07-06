package appwindow

import "image/color"

var (
	unitGreen             = color.RGBA{140, 240, 180, 128}
	unitDarkGreen         = color.RGBA{0, 0x44, 0x22, 0xFF}
	unitBlue              = color.RGBA{140, 180, 240, 128}
	unitDarkBlue          = color.RGBA{0, 0x22, 0x44, 0xFF}
	commandGreen          = color.RGBA{40, 240, 180, 128}
	commandDarkGreen      = color.RGBA{0, 0x88, 0x44, 0xff}
	commandBlue           = color.RGBA{40, 180, 240, 128}
	commandDarkBlue       = color.RGBA{0, 0x44, 0x88, 0xff}
	commandRed            = color.RGBA{200, 0, 0, 128}
	commandDarkRed        = color.RGBA{200, 0, 0, 255}
	mapGrid               = color.RGBA{0x44, 0x44, 0x44, 0xFF}
	mapBlue               = color.RGBA{30, 60, 200, 200}
	mapDeepBlue           = color.RGBA{10, 30, 100, 200}
	mapHillFill           = color.RGBA{0x71, 0x66, 0x60, 0xAA}
	mapHillStroke         = color.RGBA{0x32, 0x2d, 0x2a, 0xff}
	mapWoodsFill          = color.RGBA{0, 40, 0, 64}
	mapTownFill           = color.RGBA{0x5A, 0x51, 0x4C, 0xFF}
	mapTownStroke         = color.RGBA{32, 32, 32, 200}
	mapUnitFill           = color.RGBA{0, 0, 128, 200}
	mapEnemyFill          = color.RGBA{0x88, 0, 0, 200}
	mapUnitStroke         = color.RGBA{0, 0, 0, 0xFF}
	mapUnitSelectedFill   = color.RGBA{0, 0x44, 0xFF, 0xFF}
	mapUnitSelectedStroke = color.RGBA{0xFF, 0xAA, 0, 0xFF}
	denoteUnit            = color.RGBA{200, 200, 200, 0xff}
	mapUnitCanOrder       = color.RGBA{0x61, 0xf4, 0x1c, 0xff}
	mapUnitHasOrder       = color.RGBA{0, 0, 200, 0xff}
	mapUnitOrdersStroke   = color.RGBA{0xc1, 0x47, 0x16, 0x88}
	mapUnitOrdersMarch    = color.RGBA{0x6f, 0x4e, 0x3b, 0x88}
	mapUnitOrdersAttack   = color.RGBA{0xc1, 0x2e, 0x11, 0xAA}
	mapUnitOrdersEngage   = color.RGBA{0x55, 0x55, 0xff, 0xAA}
	mapUnitOrdersCharge   = color.RGBA{0x64, 0x31, 0x20, 0xAA}
	mapUnitOrdersFire     = color.RGBA{0xa6, 0x12, 0x3a, 0xAA}
	mapUnitOrdersPursuit  = color.RGBA{0x6f, 0x32, 0x9a, 0x88}
	plotBackground        = color.RGBA{0x88, 0x88, 0x88, 0xff}
)
