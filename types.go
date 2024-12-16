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
	ss := Split(string(l))
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

// Map alias type to support comma-separated key=value pairs
type Map string

func NewMap(m map[string]string) Map {
	strs := make([]string, 0, len(m))
	for k, v := range m {
		strs = append(strs, k+"="+v)
	}
	return Map(strings.Join(strs, ","))
}

func (m Map) ToMap() map[string]string {
	out := map[string]string{}
	pairs := Split(string(m))
	m.NavPairs(pairs, func(key, val string) {
		out[key] = val
	})
	return out
}

func (m Map) ToEnv() []string {
	pairs := Split(string(m))
	out := make([]string, 0, len(pairs))
	m.NavPairs(pairs, func(key, val string) {
		out = append(out, key+"="+val)
	})
	return out
}

func (m Map) NavPairs(pairs []string, f func(key, val string)) {
	for _, v := range pairs {
		kv := strings.Split(v, "=")
		if len(kv) != 2 {
			continue
		}
		f(kv[0], kv[1])
	}
}

func (m Map) String() string {
	return string(m)
}

var Splitter = ","

func Split(s string) []string {
	return strings.Split(s, Splitter)
}
