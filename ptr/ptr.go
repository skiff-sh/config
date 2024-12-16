package ptr

func Ptr[T any](t T) *T {
	return &t
}

func OrNil[T comparable](t T) *T {
	var c T
	if t == c {
		return nil
	}
	return Ptr(t)
}

func Deref[T any](t *T) (out T) {
	if t == nil {
		return
	}
	return *t
}

func FirstNonZeroValue[T comparable](vals ...T) (out T) {
	for _, v := range vals {
		if v != out {
			return v
		}
	}
	return
}

// FirstNonZeroOrDefaultValue returns the first value from vals that does not equal def
// and is not zero. If no value can be found, def is returned.
func FirstNonZeroOrDefaultValue[T comparable](def T, vals ...T) (out T) {
	for _, v := range vals {
		if v != def && v != out {
			return v
		}
	}

	return def
}
