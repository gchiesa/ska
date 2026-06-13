package skaffolder

func mapStringToMapInterface(m map[string]string) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range m {
		res[k] = v
	}
	return res
}

func mapInterfaceToString(m map[string]interface{}) map[string]string {
	res := make(map[string]string)
	for k, v := range m {
		res[k] = v.(string)
	}
	return res
}
