package strconv

import (
	"strconv"
)

func ParseUintOrFallback(s string, base int, bitSize int, defaultValue uint64) (uint64, error) {
	ns, err := strconv.ParseUint(s, base, bitSize)
	if err != nil {
		return defaultValue, err
	}

	return ns, nil
}

func ParseInt64Batch(m map[string]string) (map[string]int64, error) {
	nm := make(map[string]int64, 0)

	for k, v := range m {
		ns, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		nm[k] = ns
	}

	return nm, nil
}
