package query

import "fmt"

type selector struct {
	name   string
	filter *filter
}

func Selector(name string, filter *filter) *selector {
	return &selector{
		name:   name,
		filter: filter,
	}
}

func (s *selector) String() string {
	out := s.name
	if s.filter != nil {
		out += fmt.Sprintf(":filter(%s)", s.filter.String())
	}
	return out
}
