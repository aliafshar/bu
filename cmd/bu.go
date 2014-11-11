package main

import (
	"github.com/aliafshar/bu"
	"github.com/aliafshar/toylog"
	"gopkg.in/alecthomas/kingpin.v1"
	"os"
)

var (
	app        = kingpin.New("bu", "A build utility.")
	bufile     = app.Flag("bufile", "Path to bu file.").Default("main.bu").Short('f').ExistingFile()
	version    = app.Flag("version", "Print the bu version and exit.").Short('v').Bool()
	targetName = app.Arg("target", "Execute the named target.").String()
  targetArgs = app.Arg("args", "Arguments to pass to the bu target.").Strings()
  versionInfo = "bu, version " + bu.BuVersion
)

func showVersion() {
	toylog.Infoln(versionInfo)
}

func main() {
	app.Parse(os.Args[1:])
	if *version {
	  showVersion()
		return
	}
  toylog.Infof(versionInfo + ", loading %q", *bufile)
  s := bu.NewScript(*bufile, *targetArgs)
	t := s.Target(*targetName)
	if t == nil {
		toylog.Fatalf("target not found %q", *targetName)
	}
	bu.Run(s, t)
}
