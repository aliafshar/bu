```
22:16 I ┏━ ┃ ┃  
22:16 I ┏━┃┃ ┃   bu, version 0.0
22:16 I ━━ ━━┛  

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
22:16 I ●(cyan) [/tmp/main.bu:demo] "echo Hello, world!"
Hello, world!
22:16 I ●(green) 0 [/tmp/main.bu:demo]

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
22:16 I ●(cyan) [/tmp/main.bu:build] "echo a dependency"
a dependency
22:16 I ●(green) 0 [/tmp/main.bu:build]
22:16 I ●(cyan) [/tmp/main.bu:demo] "echo Hello, world!"
Hello, world!
22:16 I ●(green) 0 [/tmp/main.bu:demo]

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
22:16 I ●(cyan) [/tmp/main.bu:demo] "for i in range(5):\n  print i"
0
1
2
3
4
22:16 I ●(green) 0 [/tmp/main.bu:demo]

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
22:16 I ●(cyan) [/tmp/main.bu:make] "echo Blah > my_file.txt"
22:16 I ●(green) 0 [/tmp/main.bu:make]
22:16 I ●(cyan) [/tmp/main.bu:demo] "cat my_file.txt\nrm my_file.txt"
Blah
22:16 I ●(green) 0 [/tmp/main.bu:demo]

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
22:16 I ●(cyan) [/tmp/main.bu:demo] "echo piped\necho banana" | "wc -c" | "wcalc -h"
 = 0xd
22:16 I ●(green) 0 | 0 | 0 [/tmp/main.bu:demo]

```

Here the output of the `pipe` target is piped into the count target and then the
hex target. Of course, all dependencies will be first run.

# Watches

Targets can be restarted based on watching a file for modification. This is
probably only useful for long-running targets.

```bu
a_dep:
  echo hello

demo: a_dep ^example.bu
  sleep 5
```

```bu-out
