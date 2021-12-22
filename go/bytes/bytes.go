package bytes

func Truncate(s []byte, n int) []byte {
	if n < 0 {
		n = 0
	}

	if len(s) <= n {
		return s
	}

	return s[:n]
}
