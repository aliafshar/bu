package bu

type target struct {
	name     string
	body     string
	shell    string
	deps     []dependency
	pipe     []dependency
	redirect *redirect
	watch    string
}

type result struct {
	err  error
	desc string
}

func (r *result) success() bool {
	return r.err == nil
}
