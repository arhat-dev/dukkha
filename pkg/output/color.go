package output

import "github.com/fatih/color"

func PickColor(i int) (prefixColor, outputColor *color.Color) {
	// TODO: generate colors dynamically with consistent result

	var colorList = [][2]*color.Color{
		{color.New(color.FgHiWhite), color.New(color.FgWhite)},
		{color.New(color.FgHiCyan), color.New(color.FgCyan)},
		{color.New(color.FgHiGreen), color.New(color.FgGreen)},
		{color.New(color.FgHiMagenta), color.New(color.FgMagenta)},
		{color.New(color.FgHiYellow), color.New(color.FgYellow)},
		{color.New(color.FgHiBlue), color.New(color.FgBlue)},
		{color.New(color.FgHiRed), color.New(color.FgRed)},
	}

	colorSet := colorList[i%len(colorList)]

	return colorSet[0], colorSet[1]
}
