# 𝄽 bu

Bu is a tool to help you run common tasks. It is something like a simple version
of make. You define a set of tasks and it will run them.

    run:
        echo I would run something.

When running `bu run`,

    I: bu, version 0.0, loading "Bufile"
    I: < "echo I run something.". [worker:0]
    I: > I run something.

## Imports

    < foo.bu

Will import foo.bu from the system path, which defaults to resolving, in order:

* Current working directory `.`
* Bu home directory `~/.bu`

Comments  
Line comments only. Non-line comments are undefined.
