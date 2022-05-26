package clihelper

import "github.com/spf13/pflag"

func InitFlagSet(fs *pflag.FlagSet, name string) {
	fs.Init(name, pflag.ContinueOnError)
	fs.SetInterspersed(true)
	fs.SortFlags = true
}
