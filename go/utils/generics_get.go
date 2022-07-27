package utils

func GetValueOrFallback[T comparable] (v, defaultValue T) T {
  var t T
  if v == t  {
	return defaultValue
  }

  return v
}

func Pointer[T any](v T) *T {  
  return &v
}
