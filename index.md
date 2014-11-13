Bu is a tool to help you run common tasks. It is something like a simple version
of GNU make with some additional features. You define a set of tasks and it will
run them. It features: **targets with dependencies**, **script imports**,
**using Bash or Python**, **task parallelism**, **command line inputs**,
and **variables**.

Here is a tiny example.

```bu
demo:
  echo Hello, world!
```

```bu-out
I: bu, version 0.0, loading "/tmp/tmpbgUyzZ"
I: > [demo] sh:"echo Hello, world!"
Hello, world!
I: < [demo] success
```

The target is executed with:

    $ bu run

And we get the following output:

    I: bu, version 0.0, loading "Bufile"
    I: < "echo I run something.". [worker:0]
    I: > I run something.

## Targets

    <target name>: [dependencies...] [!type] 
        <script body>

for example,

```bu
demo: build
  echo Hello, world!

build:
  echo a dependency
```

```bu-out
I: bu, version 0.0, loading "/tmp/tmpwGotJS"
I: > [build] sh:"echo a dependency"
a dependency
I: < [build] success
I: > [demo] sh:"echo Hello, world!"
Hello, world!
I: < [demo] success
```

is a target named `run` that depends on a target named `build` that runs the
shell command `go run cmd/bu.go`.

### Target types

Currently only shell and python are supported. Shell is the default, so no type
is required to be passed explicitly. For a Python target, add the type.

```bu
demo: !py
  for i in range(5):
    print i
```

```bu-out
I: bu, version 0.0, loading "/tmp/tmpXCjuMF"
I: > [demo] py:"for i in range(5):\n  print i"
0
1
2
3
4
I: < [demo] success
```

### Indentation

Target bodies must be indented by any whitespace, tab or space. Indentation must
be consistent for Python scripts since Python is sensitive to this.

## Variables


Single line variables are defined with the `=` operator, like so:

```bu
DEMO =I am the variable content

demo:
  echo $DEMO
```

```bu-out
I: bu, version 0.0, loading "/tmp/tmp5_JGHe"
I: > [demo] sh:"echo $DEMO"
I am the variable content
I: < [demo] success
```

Multiline variables are defined with the `=|` operator followed by a block.

```bu
DEMO =|
    I
    am
    the variable
    content

demo:
  echo "$DEMO"
```

```bu-out
I: bu, version 0.0, loading "/tmp/tmpOB4Ttn"
I: > [demo] sh:"echo \"$DEMO\""
I
am
the variable
content
I: < [demo] success
```

Defines a variable `myvariable`. Quoting is not required as the variable value
is taken to the end of the line.

Variables are injected into the environment,
where they can be used directly in targets as `$myvariable` in shell, or
`os.getenv('myvariable')` from Python.

Note: Because variables are injected only into the environment, they will not be
used in target names and dependencies.

## Positional arguments 

`$0`, `$1`, etc (in shell) and `sys.args` (in Python) are available as
additional arguments passed to the bu invocation. Consider this target:

```bu
demo:
  echo Hi, "$0"
```

```bu-out
I: bu, version 0.0, loading "/tmp/tmpmDTwM6"
I: > [demo] sh:"echo Hi, \"$0\""
Hi, demo
I: < [demo] success
```

and this invocation:

    $ bu demo FirstArgument

## Questions

```bu
demo ? n
  Are you sure? (y/n)

danger: demo
  if [ $demo -eq y ]; then
    echo Confirmed, continuing
  fi
```

```bu-out
I: bu, version 0.0, loading "/tmp/tmpYqEg7U"
I: > [demo] question: "Are you sure? (y/n)"
[1mAre you sure? (y/n)[0m (default=[1m[34mn[0m[0m) > E: < [demo] failure, EOF
I: < [demo] success $demo=""
```

Will prompt the user on the command line and store the value in the variable
`confirm` with a default value of `n`. Questions are targets and can be depended
on by other targets.

Default values are optional, with the syntax:

    <name> ? [default]
        <question>

## Imports

    < foo.bu

Will import foo.bu from the system path, which defaults to resolving, in order:

* Current working directory `.`
* Bu home directory `~/.bu`

## Comments  

Line comments only. Non-line comments are undefined, especially in situations
where values are taken to the end of a line, e.g. variable definitions

    myvariable = I am the value # this comment will be part of the value

## Differenced from GNU make

* Each target is executed in the same shell
* File existence is not explicitly taken to imply a dependency satisfaction
