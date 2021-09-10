package output

import (
	"arhat.dev/dukkha/pkg/dukkha"
	"github.com/muesli/termenv"
)

func PickColor(i int) (prefixColor, outputColor dukkha.TermColor) {
	// TODO: generate colors dynamically with consistent result

	var colorList = [][2]termenv.ANSIColor{
		{termenv.ANSIBrightWhite, termenv.ANSIWhite},
		{termenv.ANSIBrightCyan, termenv.ANSICyan},
		{termenv.ANSIBrightGreen, termenv.ANSIGreen},
		{termenv.ANSIBrightMagenta, termenv.ANSIMagenta},
		{termenv.ANSIBrightYellow, termenv.ANSIYellow},
		{termenv.ANSIBrightBlue, termenv.ANSIBlue},
		{termenv.ANSIBrightRed, termenv.ANSIRed},
	}

	colorSet := colorList[i%len(colorList)]

	return dukkha.TermColor(colorSet[0]), dukkha.TermColor(colorSet[1])
}
