package os

func GetValueOrFallback[T comparable](v, defaultValue T) T {
	var zeroE T
	if v != zeroE {
		return v
	}

	return defaultValue
}
