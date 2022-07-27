package utils

func GetValueOrFallback[T comparable] (v, defaultValue T) T {
  var t T
  if v == t  {
	return defaultValue
  }

  return v
}
