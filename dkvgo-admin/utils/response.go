package utils

func ErrorMap(message string) map[string]interface{} {
	errMap := make(map[string]interface{})
	errMap["success"] = false
	errMap["message"] = message
	return errMap
}

func SuccessMap() map[string]interface{} {
	sucMap := make(map[string]interface{})
	sucMap["success"] = true
	sucMap["message"] = "SUCCESS"
	return sucMap
}