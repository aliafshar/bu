
# 𝄽 bu

Bu is a tool to help you run common tasks. It is something like a simple version
of make. You define a set of tasks and it will run them.

- [Targets](#targets)
  - [Target types](#target-types)
  - [Indentation](#indentation)
- [Variables](#variables)
- [Imports](#imports)
- [Comments](#comments)

Here is a tiny example.

    run:
        echo I run something.

When running `bu run`,

    I: bu, version 0.0, loading "Bufile"
    I: < "echo I run something.". [worker:0]
    I: > I run something.

## Targets

    <target name>: [dependencies...] [!type] 
        <script body>

for example,

    run: build
        go run cmd/bu.go

is a target named `run` that depends on a target named `build` that runs the
shell command `go run cmd/bu.go`.

### Target types

Currently only shell and python are supported. Shell is the default, so no type
is required to be passed explicitly. For a Python target, add the type.

    run: !py
        for i in range(5):
            print i

### Indentation

Target bodies must be indented by any whitespace, tab or space. Indentation must
be consistent for Python scripts since Python is sensitive to this.

## Variables

    myvariable = I am the variable content

Defines a variable `myvariable`. Quoting is not required as the variable value
is taken to the end of the line.

Variables are injected into the environment,
where they can be used directly in targets as `$myvariable` in shell, or
`os.getenv('myvariable')` from Python.

Note: Because variables are injected only into the environment, they will not be
used in target names and dependencies.

## Imports

    < foo.bu

Will import foo.bu from the system path, which defaults to resolving, in order:

* Current working directory `.`
* Bu home directory `~/.bu`

## Comments  

Line comments only. Non-line comments are undefined, especially in situations
where values are taken to the end of a line, e.g. variable definitions

    myvariable = I am the value # this comment will be part of the value
