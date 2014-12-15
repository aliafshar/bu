package bu

import (
	"testing"
  "strings"
)

var (
  b1 = `# Bu example`
  b2 = `my_var =42`
  b3 = `'quoted something'`
)

func lexSnippet(s string) []*token {
  o := make(chan *token)
  var ts []*token
  go lex(strings.NewReader(s), o)
  for {
    t := <-o
    ts = append(ts, t)
    if t.is(tokenEof) {
      break
    }
  }
  return ts
}

func Test_Ends(t *testing.T) {
  ts := lexSnippet(b1)
  if !ts[0].is(tokenSof) {
    t.Fatal("first token is not SOF")
  }
  if !ts[1].is(tokenSol) {
    t.Fatal("second token is not SOL")
  }
  if !ts[len(ts) - 1].is(tokenEof) {
    t.Fatal("last token is not EOF")
  }
  if !ts[len(ts) - 2].is(tokenEol) {
    t.Fatal("penultimate token is not EOL")
  }
}

func Test_Comment(t *testing.T) {
  ts := lexSnippet(b1)
  if !ts[2].is(tokenComment) {
    t.Fatal(ts)
  }
  if !ts[3].is(tokenRaw) {
    t.Fatal(ts)
  }
  if !(ts[3].value() == " Bu example") {
    t.Fatal(ts)
  }
}

func Test_Quotes(t *testing.T) {
  ts := lexSnippet(b3)
  if !ts[2].is(tokenQuote) {
    t.Fatal(ts)
  }
  if !ts[3].is(tokenName) {
    t.Fatal(ts)
  }
  if !(ts[3].value() == "quoted something") {
    t.Fatal(ts)
  }
  if !ts[4].is(tokenQuote) {
    t.Fatal(ts)
  }
}

