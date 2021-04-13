package query

import "fmt"

func Eq(key, value string) *filter {
	return &filter{
		name:  "eq",
		key:   key,
		value: value,
	}
}

func And(f1, f2 *filter) *filter {
	return &filter{
		name:  "and",
		key:   f1.String(),
		value: f2.String(),
	}
}

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
