package config

import "strings"

// List a comma-separated list of strings e.g. val1,val2
type List string

// NewList instantiate a new List.
func NewList(ss ...string) List {
	return List(strings.Join(ss, ","))
}

// ToSlice convert the list to a slice of strings.
func (l List) ToSlice() []string {
	ss := strings.Split(string(l), ",")
	out := make([]string, 0, len(ss))
	for _, v := range ss {
		if v == "" {
			continue
		}
		out = append(out, v)
	}
	return out
}

func (l List) String() string {
	return string(l)
}
