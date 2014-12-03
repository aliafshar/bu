package bu

import(
  "io"
)

type targetType func(*node) target

type target interface {
	Name() string
	Deps() []dependency
	Run(*runContext) result
	Desc() string
}

type result interface {
	Success() bool
	Err() error
	Desc() string
}

type runContext struct {
  in io.Reader
  out io.Writer
  worker *worker
  script *script
}

