package util

// MapToCmdArgs convert a map[string]string to args string slice
func MapToCmdArgs(m map[string]string, tag string, fields ...string) []string {
	var args = make([]string, 0, len(m)*2)
	for _, field := range fields {
		if value, ok := m[field]; ok {
			args = append(args, tag+field, value)
		}
	}
	return args
}
