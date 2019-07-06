package mapeditor

import "image/color"

var (
	mapGrid       = color.RGBA{0x44, 0x44, 0x44, 0xFF}
	mapSelect     = color.RGBA{0xff, 0, 0, 0xff}
	mapBlue       = color.RGBA{30, 60, 200, 200}
	mapDeepBlue   = color.RGBA{10, 30, 100, 200}
	mapHillFill   = color.RGBA{0x71, 0x66, 0x60, 0xAA}
	mapHillStroke = color.RGBA{0x32, 0x2d, 0x2a, 0xff}
	mapWoodsFill  = color.RGBA{0, 40, 0, 64}
	mapTownFill   = color.RGBA{0x5A, 0x51, 0x4C, 0xFF}
	mapTownStroke = color.RGBA{32, 32, 32, 200}
)
