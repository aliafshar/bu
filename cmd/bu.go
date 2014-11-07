package main

import (
	"github.com/aliafshar/bu"
	"github.com/aliafshar/toylog"
	"gopkg.in/alecthomas/kingpin.v1"
	"os"
)

var (
	app        = kingpin.New("bu", "A build utility.")
	bufile     = app.Flag("bufile", "Path to bu file.").Default("main.bu").Short('f').String()
	version    = app.Flag("version", "Print the bu version and exit.").Short('v').Bool()
	targetName = app.Arg("target", "Execute the named target.").String()
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
	//s, err := bu.Load(*bufile)
	//if err != nil {
	//	toylog.Fatalf("failed to load bu file (%v)", err)
	//}
	//t := s.Target(*targetName)
  s := bu.NewScript(*bufile)
	t := s.Target(*targetName)
	if t == nil {
		toylog.Fatalf("target not found %q", *targetName)
	}
	bu.Run(s, t)
}
