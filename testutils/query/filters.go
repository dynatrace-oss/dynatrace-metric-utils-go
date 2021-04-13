package query

import "fmt"

// Eq returns a filter which checks that the dimension equals the given value
func Eq(key, value string) *filter {
	return &filter{
		name:  "eq",
		key:   key,
		value: value,
	}
}

// And takes two filters and returns a filter which ensures they are both true
func And(f1, f2 *filter) *filter {
	return &filter{
		name:  "and",
		key:   f1.String(),
		value: f2.String(),
	}
}

// Or takes two filters and returns a filter which ensures at least one is true
func Or(f1, f2 *filter) *filter {
	return &filter{
		name:  "or",
		key:   f1.String(),
		value: f2.String(),
	}
}

type filter struct {
	name  string
	key   string
	value string
}

func (f *filter) String() string {
	return fmt.Sprintf("%s(%s,%s)", f.name, f.key, f.value)
}
