package main

import (
	"bitbucket.org/kardianos/osext"
	"github.com/aliafshar/bu"
	"github.com/aliafshar/toylog"
	"gopkg.in/alecthomas/kingpin.v1"
  "io"
	"os"
  "strings"
)

var (
	app         = kingpin.New("bu", "A build utility.")
	bufile      = app.Flag("bufile", "Path to bu file.").Default("main.bu").Short('f').ExistingFile()
	version     = app.Flag("version", "Print the bu version and exit.").Short('v').Bool()
	debug       = app.Flag("debug", "Verbose logging.").Short('d').Bool()
	quiet       = app.Flag("quiet", "Don't be so noisy.").Short('q').Bool()
	content     = app.Flag("content", "The script contents. Overrides bufile.").Short('c').String()
	targetName  = app.Arg("target", "Execute the named target.").String()
	targetArgs  = app.Arg("args", "Arguments to pass to the bu target.").Strings()
	versionInfo = "bu, version " + bu.BuVersion
)

var (
	logo1 = "┏━ ┃ ┃  "
	logo2 = "┏━┃┃ ┃  "
	logo3 = "━━ ━━┛  "
)

func mustExecutable() interface{} {
	filename, err := osext.Executable()
	if err != nil {
		return err
	}
	return filename
}


func main() {
	app.Parse(os.Args[1:])
	if *debug {
		toylog.Verbose()
	}
	toylog.Debugf("starting from %q\n", mustExecutable())
	if !*quiet {
		toylog.Infoln(logo1)
		toylog.Infoln(logo2, versionInfo)
		toylog.Infoln(logo3)
	} else {
		toylog.Debugln("quiet mode")
	}
	if *version {
		toylog.Debugln("version and exit")
		return
	}

  var r io.Reader
  var n string
  if *content == "" {
    f, err := os.Open(*bufile)
	  if err != nil {
      toylog.Fatalln("Unable to open bu file.", *bufile, err)
	  }
	  defer f.Close()
    r = f
    n = *bufile
  } else {
    r = strings.NewReader(*content)
    n = "<command line>"
  }
	toylog.Debugf("running %q from %q with %q\n", *targetName, n, *targetArgs)
	bu.Run(r, n, *targetName, *targetArgs...)
}
