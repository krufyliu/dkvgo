package utils

func ErrorMap(message string) map[string]interface {
	var errMap = make(map[string]interface)
	errMap["success"] = false
	errMap["message"] = message
	return errMap
}