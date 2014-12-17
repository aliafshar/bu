package bu

import (
	"io"
	"os"

	"github.com/aliafshar/toylog"
)

type redirect struct {
	ifile string
	ofile string
}

func (r *redirect) out() io.Writer {
	if r.ofile == "" {
		return os.Stdout
	}
	f, err := os.Create(r.ofile)
	if err != nil {
		toylog.Errorln("Unable to open output file", r.ofile, err)
		return os.Stdout
	}
	return f
}

func (r *redirect) in() io.Reader {
	if r.ifile == "" {
		return nil
	}
	f, err := os.Open(r.ifile)
	if err != nil {
		toylog.Errorln("Unable to open input file", r.ifile, err)
	}
	return f
}
