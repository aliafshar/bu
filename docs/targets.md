---
draft: true
title: "Targets"
description: Details on what targets are and how to use them
menu:
  docs:
    weight: -100
date: 2014-01-01
---

Targets are the fundamental building block of bu. The are a group of dependent
tasks with multiple options. The syntax is fundamentally similar to how make
does things, and many make scripts will run in bu.

# Your first target

Let's look at the simplest
target.

```bu
demo:
  echo hello bu world 𝄽
```

As you can see bu executes the shell script and displays the output. A target definition looks like this.

```bu-spec
<target name>: [target dependencies...] [?file dependencies] [!type] [>outfile] [<infile] [|pipe] [^watch]
    <script body>
```

Note that the body must be indented, but unlike Make, tabs or spaces are fine.


# Dependencies

All dependencies of a target must be met for it to run.

## Target dependencies

A target can depend on another target by name.

```bu
t1:
  echo I am depended on by demo

demo: t1
  echo I depend on t1
```

## File Dependencies

A target may explicitly depend on the existence of a file or directory.

```bu
make:
  echo Blah > my_file.txt

demo: make ?my_file.txt
  cat my_file.txt
  rm my_file.txt
```

## Web dependencies

A target can depend on a web page being present and returning 200. The page will
be polled until it can be contacted.

```bu
demo: @example.com
  echo example.com is up
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

Here the output of the `pipe` target is piped into the count target and then the
hex target. Of course, all dependencies will be first run.

# Watches

Targets can be restarted based on watching a file for modification. This is
probably only useful for long-running targets.

```bu
a_dep:
  echo hello

demo: a_dep ^example.bu
  sleep 0
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

Similarly a file can be used for input on standard input.

```bu
make: >my_file.txt
  print "Save me in a file"

demo: <my_file.txt !py
  import sys
  print sys.stdin.read()
```

# Shell types

Currently only shell and python are supported. Shell (bash) is the default, so no type
is required to be passed explicitly. For a Python target, add the type as `!py`,
look:

```bu
demo: !py
  print "hello bu world 𝄽"
```


# Indentation and whitespace

Target bodies must be indented by any whitespace, tab or space. Indentation must
be consistent for Python scripts since Python is sensitive to this.

```bu
demo: !py
  for i in range(5):
  print i
```
