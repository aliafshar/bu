package main

import (
	"bytes"
	"github.com/aliafshar/toylog"
	"gopkg.in/alecthomas/kingpin.v1"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
)

var (
	app  = kingpin.New("bu-replace", "Find bu snippets and run them.")
	in   = app.Flag("in", "Input file").Required().String()
	out  = app.Flag("out", "Output file").Required().String()
	expr = regexp.MustCompile("(?s)```bu\\n(.+?)```")
)

func runSnippet(s string) string {
	cmd := exec.Command("bu", "-q", "-c", s, "demo")
	bs, err := cmd.CombinedOutput()
	if err != nil {
		toylog.Fatalln(err)
	}
	return string(bs)
}

func replacer(s string) string {
	b := bytes.NewBufferString(s)
	b.WriteString("\n\n```bu-output\n")
	bd := expr.FindString(s)
	b.WriteString(runSnippet(bd))
	b.WriteString("```\n\n")
	return b.String()
}

func replaceFile(in, out string) {
	bs, err := ioutil.ReadFile(in)
	if err != nil {
		toylog.Fatalln(err)
	}
	o := expr.ReplaceAllStringFunc(string(bs), replacer)
	ioutil.WriteFile(out, []byte(o), 0644)
}

func main() {
	app.Parse(os.Args[1:])
	replaceFile(*in, *out)
}
