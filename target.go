package bu

type targetType func(*node) target

type target interface {
	Name() string
	Deps() []dependency
	Run() result
  Desc() string
}

type result interface {
  Success() bool
  Err() error
  Desc() string
}
