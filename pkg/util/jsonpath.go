package util

import (
	"encoding/json"
	"k8s.io/client-go/util/jsonpath"
	"strings"
)

func QueryJSONString(jsonText, jsonPathQuery string) (string, error) {
	var jsonData interface{}
	if err := json.Unmarshal([]byte(jsonText), &jsonData); err != nil {
		return "", err
	}
	jp := jsonpath.New("parser")
	if err := jp.Parse(jsonPathQuery); err != nil {
		return "", err
	}
	var buff strings.Builder
	if err := jp.Execute(&buff, jsonData); err != nil {
		return "", err
	}
	return buff.String(), nil
}
