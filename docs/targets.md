---
draft: true
title: "Targets"
description: Details on what targets are and how to use them
menu:
  docs:
    weight: -100
date: 2014-01-01
---

# Targets

Targets are the fundamental building block of bu. The are a group of dependent
tasks with multiple options. The syntax is fundamentally similar to how make
does things, and many make scripts will run in bu.


## Your first target

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

# Choosing a shell

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

