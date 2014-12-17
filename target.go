package bu

type target struct {
	name     string
	body     string
	shell    string
	deps     []dependency
	pipe     []dependency
	redirect *redirect
}

type result struct {
	err  error
	desc string
}

func (r *result) success() bool {
	return r.err == nil
}

func combinedResult(rs []*result) *result {
	r := &result{}
	for _, rr := range rs {
		if rr.err != nil {
			r.err = rr.err
			r.desc = r.desc + " 1"
		} else {
			r.desc = r.desc + " 0"
		}
	}
	return r
}
