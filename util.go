package dkvgo

// mapToStringSlice convert a map[string]string to args string slice
func MapToCmdArgs(m map[string]string, tag string) []string {
	var args = make([]string, 0, len(m)*2)
	for key, value := range m {
		args = append(args, tag+key, value)
	}
	return args
}
