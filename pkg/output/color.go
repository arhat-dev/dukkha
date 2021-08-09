package output

import (
	"github.com/muesli/termenv"
)

func PickColor(i int) (prefixColor, outputColor termenv.Color) {
	// TODO: generate colors dynamically with consistent result

	var colorList = [][2]termenv.Color{
		{termenv.ANSIBrightWhite, termenv.ANSIWhite},
		{termenv.ANSIBrightCyan, termenv.ANSICyan},
		{termenv.ANSIBrightGreen, termenv.ANSIGreen},
		{termenv.ANSIBrightMagenta, termenv.ANSIMagenta},
		{termenv.ANSIBrightYellow, termenv.ANSIYellow},
		{termenv.ANSIBrightBlue, termenv.ANSIBlue},
		{termenv.ANSIBrightRed, termenv.ANSIRed},
	}

	colorSet := colorList[i%len(colorList)]

	return colorSet[0], colorSet[1]
}
