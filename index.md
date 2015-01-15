```
15:31 I ┏━ ┃ ┃  
15:31 I ┏━┃┃ ┃   bu, version 0.0
15:31 I ━━ ━━┛  

```

Note: the full content of this document, with executed snippets is available at
[in the documentation](http://aliafshar.github.io/bu)

Bu is a tool to help you run common tasks. It is something like Gulp that looks
like Make. You define a set of tasks and it will run them. It features:
**targets with dependencies**, **script imports**, **using Bash or Python**,
**task parallelism**, **command line inputs**, and **variables**.

Here is a tiny example.

```bu
demo:
  echo Hello, world!
```

```bu-out
15:31 I ●(cyan) [/tmp/tmpDC6ZXz:demo] "echo Hello, world!"
Hello, world!
15:31 I ●(green) 0 [/tmp/tmpDC6ZXz:demo]

```

The target is executed with:

    $ bu demo

Or since it is the first target, simply:

    $ bu

# Usage

```
usage: bu [<flags>] [<target> [<args>]]

A build utility.

Flags:
  --help         Show help.
  -f, --bufile=main.bu  Path to bu file.
  -v, --version  Print the bu version and exit.
  -d, --debug    Verbose logging.
  -l, --list     List targets.
  -q, --quiet    Don't be so noisy.

Args:
  [<target>]  Execute the named target.
  [<args>]    Arguments to pass to the bu target.

```

# Targets

Targets are the unit of work. They support a number of options.

```bu-spec
<target name>: [target dependencies...] [?file dependencies] [!type] [>outfile] [<infile] [|pipe]
  <script body>
```

for example,

```bu
demo: build
  echo Hello, world!

build:
  echo a dependency
```

```bu-out
15:31 I ●(cyan) [/tmp/tmpz9HfLn:build] "echo a dependency"
a dependency
15:31 I ●(green) 0 [/tmp/tmpz9HfLn:build]
15:31 I ●(cyan) [/tmp/tmpz9HfLn:demo] "echo Hello, world!"
Hello, world!
15:31 I ●(green) 0 [/tmp/tmpz9HfLn:demo]

```

is a target named `run` that depends on a target named `build` that runs the
shell command `go run cmd/bu.go`.

## Target types

Currently only shell and python are supported. Shell is the default, so no type
is required to be passed explicitly. For a Python target, add the type.

## Indentation

Target bodies must be indented by any whitespace, tab or space. Indentation must
be consistent for Python scripts since Python is sensitive to this.

```bu
demo: !py
  for i in range(5):
    print i
```

```bu-out
15:31 I ●(cyan) [/tmp/tmpbDUxnp:demo] "for i in range(5):\n  print i"
0
1
2
3
4
15:31 I ●(green) 0 [/tmp/tmpbDUxnp:demo]

```

# Dependencies

## File Dependencies

A target may explicitly depend on the existence of a file or directory.

```bu
make:
  echo Blah > my_file.txt

demo: make ?my_file.txt
  cat my_file.txt
  rm my_file.txt
```

```bu-out
15:31 I ●(cyan) [/tmp/tmp58oTnV:make] "echo Blah > my_file.txt"
15:31 I ●(green) 0 [/tmp/tmp58oTnV:make]
15:31 I ●(cyan) [/tmp/tmp58oTnV:demo] "cat my_file.txt\nrm my_file.txt"
Blah
15:31 I ●(green) 0 [/tmp/tmp58oTnV:demo]

```

# Pipes

Targets can be piped into eachother.

```bu
count:
  wc -c

hex:
  wcalc -h

demo: | count | hex
  echo piped
  echo banana
```

```bu-out
15:31 I ●(cyan) [/tmp/tmpKiPs7g:demo] "echo piped\necho banana" | "wc -c" | "wcalc -h"
 = 0xd
15:31 I ●(green) 0 | 0 | 0 [/tmp/tmpKiPs7g:demo]

```

Here the output of the `pipe` target is piped into the count target and then the
hex target. Of course, all dependencies will be first run.

# Watches

Targets can be restarted based on watching a file for modification. This is
probably only useful for long-running targets.

```
a_dep:
  echo hello

demo: a_dep ^example.bu
  sleep 5
```

Will restart the `watch` target every time `example.bu` file is modified. It
will handle stopping the running process.

# Redirects

Target output can be redirected to a file. This is useful when using shells that
don't have redirection, like Python.

```bu
demo: >my_file.txt !py
  print "Save me in a file"
```

```bu-out
15:31 I ●(cyan) [/tmp/tmpDKa_kD:demo] "print \"Save me in a file\""
15:31 I ●(green) 0 [/tmp/tmpDKa_kD:demo]

```

Similarly a file can be used for input on standard input.

```bu
make: >my_file.txt
  print "Save me in a file"

demo: <my_file.txt !py
  import sys
  print sys.stdin.read()
```

```bu-out
15:31 I ●(cyan) [/tmp/tmpQEUT7_:demo] "import sys\nprint sys.stdin.read()"
Save me in a file

15:31 I ●(green) 0 [/tmp/tmpQEUT7_:demo]

```


# Variables

Bu does not have variables of it's own. Only environment variables are supported.
These are passed to all targets.

Single line variables are defined with the `=` operator, like so:

```bu
DEMO =I am the variable content

demo:
  echo $DEMO
```

```bu-out
15:31 I ●(cyan) [/tmp/tmpO2mt4O:demo] "echo $DEMO"
I am the variable content
15:31 I ●(green) 0 [/tmp/tmpO2mt4O:demo]

```

Multiline variables are defined exactly the same way in an indented block.

```bu
DEMO =
    I
    am
    the variable
    content

demo:
  echo "$DEMO"
```

```bu-out
15:31 I ●(cyan) [/tmp/tmpoeXNbV:demo] "echo \"$DEMO\""
I
am
the variable
content
15:31 I ●(green) 0 [/tmp/tmpoeXNbV:demo]

```

Defines a variable `myvariable`. Quoting is not required as the variable value
is taken to the end of the line.

Variables are injected into the environment,
where they can be used directly in targets as `$myvariable` in shell, or
`os.getenv('myvariable')` from Python.

Note: Because variables are injected only into the environment, they will not be
used in target names and dependencies.

## Positional arguments 

`$0`, `$1`, `$*`, `$@` etc (in shell) and `sys.args` (in Python) are available as
additional arguments passed to the bu invocation. Consider this target:

```bu
demo:
  echo Hi, "$0"
```

```bu-out
15:31 I ●(cyan) [/tmp/tmpmoUUPh:demo] "echo Hi, \"$0\""
Hi, demo
15:31 I ●(green) 0 [/tmp/tmpmoUUPh:demo]

```

and this invocation:

    $ bu demo FirstArgument


# Imports

```bu-spec
< <filepath>
```

    < foo.bu

Will import foo.bu from the system path, which defaults to resolving, in order:

* Current working directory `.`
* Bu home directory `~/.bu`

# Comments  

Line comments only. Non-line comments are undefined, especially in situations
where values are taken to the end of a line, e.g. variable definitions

    myvariable = I am the value # this comment will be part of the value

## Differences from GNU make

* Each target is executed in the same shell
* File existence is not explicitly taken to imply a dependency satisfaction

# Prologue

To you, designers and engineers of build systems, I present Bu.

For years I have suffered your hideous constructs, your crushing assumptions,
and your bizarre choices. The only revenge left to me is to build one that is
worse, more disgusting, and infinitely heavier with dripping abomination, so
that you may suffer the pain that I have suffered. Feel the gut-wrenching taste
of bile in my mouth. Feel it, and let it sour the flavor of your day, your week
and your month. May it lay barren the fields of your productivity.

**To you, designers and engineers of build systems, I present Bu!**
