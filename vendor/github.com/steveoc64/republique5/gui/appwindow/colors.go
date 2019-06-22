package appwindow

import "image/color"

var (
	unit_green               = color.RGBA{140, 240, 180, 128}
	unit_blue                = color.RGBA{140, 180, 240, 128}
	command_green            = color.RGBA{40, 240, 180, 128}
	command_blue             = color.RGBA{40, 180, 240, 128}
	command_red              = color.RGBA{200, 0, 0, 128}
	map_grid                 = color.RGBA{0x44, 0x44, 0x44, 0xFF}
	map_grid_minor           = color.RGBA{0x44, 0x44, 0x44, 0xFF}
	map_blue                 = color.RGBA{30, 60, 200, 200}
	map_deep_blue            = color.RGBA{10, 30, 100, 200}
	map_hill_fill            = color.RGBA{0x71, 0x66, 0x60, 0xAA}
	map_hill_stroke          = color.RGBA{0x32, 0x2d, 0x2a, 0xff}
	map_woods_fill           = color.RGBA{0, 40, 0, 64}
	map_woods_stroke         = color.RGBA{40, 20, 20, 0xFF}
	map_town_fill            = color.RGBA{0x5A, 0x51, 0x4C, 0xFF}
	map_town_stroke          = color.RGBA{32, 32, 32, 200}
	map_unit_fill            = color.RGBA{0, 0, 128, 200}
	map_unit_stroke          = color.RGBA{0, 0, 0, 0xFF}
	map_unit_selected_fill   = color.RGBA{0, 0x44, 0xFF, 0xFF}
	map_unit_selected_stroke = color.RGBA{0xFF, 0xAA, 0, 0xFF}
	denote_unit              = color.RGBA{200, 200, 200, 0xff}
)