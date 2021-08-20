package strings

func GetStringOrFallback(values []string, defaultValue string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}

	return defaultValue
}
