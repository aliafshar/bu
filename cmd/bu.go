package main

import (
	"fmt"
	"github.com/aliafshar/bu"
	"github.com/aliafshar/toylog"
	"gopkg.in/alecthomas/kingpin.v1"
	"os"
)

var (
	app         = kingpin.New("bu", "A build utility.")
	bufile      = app.Flag("bufile", "Path to bu file.").Default("main.bu").Short('f').ExistingFile()
	version     = app.Flag("version", "Print the bu version and exit.").Short('v').Bool()
	debug       = app.Flag("debug", "Verbose logging.").Short('d').Bool()
	list        = app.Flag("list", "List targets.").Short('l').Bool()
	quiet       = app.Flag("quiet", "Don't be so noisy.").Short('q').Bool()
	targetName  = app.Arg("target", "Execute the named target.").String()
	targetArgs  = app.Arg("args", "Arguments to pass to the bu target.").Strings()
	versionInfo = "bu, version " + bu.BuVersion
)

var (
	logo1 = "┏━ ┃ ┃  "
	logo2 = "┏━┃┃ ┃  "
	logo3 = "━━ ━━┛  "
)

func showVersion() {
	toylog.Infoln(versionInfo)
}

func main() {
	app.Parse(os.Args[1:])
	if !*quiet {
		toylog.Infoln(logo1)
		toylog.Infoln(logo2, versionInfo)
		toylog.Infoln(logo3)
	}
	if *debug {
		toylog.Verbose()
	}
	if *version {
		return
	}
	if *list {
		for _, t := range bu.List(*bufile) {
			fmt.Println(t)
		}
		return
	}
	bu.Run(*bufile, *targetName, *targetArgs...)
}
