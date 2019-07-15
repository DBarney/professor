package storage

// Collection holds a bunch of possibly unrelated entries
// it has some convience methods attached to it.
type Collection []interface{}

func (c Collection) Each(f func(interface{})) {
	for _, e := range c {
		f(e)
	}
}

func (c Collection) First(f func(interface{}) bool) interface{} {
	for _, e := range c {
		if f(e) {
			return e
		}
	}
	return nil
}

func (c Collection) Last(f func(interface{}) bool) interface{} {
	for i := len(c) - 1; i >= 0; i-- {
		e := c[i]
		if f(e) {
			return e
		}
	}
	return nil
}

func (c Collection) Take(f func(interface{}) bool) Collection {
	res := Collection{}
	for _, e := range c {
		if f(e) {
			res = append(res, e)
		}
	}
	return res
}

func (c Collection) Toss(f func(interface{}) bool) Collection {
	res := Collection{}
	for _, e := range c {
		if !f(e) {
			res = append(res, e)
		}
	}
	return res
}
