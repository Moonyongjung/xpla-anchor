package util

import (
	"encoding/json"
	"strings"
)

func JsonUnmarshalData(jsonStruct interface{}, byteValue []byte) interface{} {
	json.Unmarshal(byteValue, &jsonStruct)

	return jsonStruct
}

func JsonMarshalData(jsonData interface{}) ([]byte, error) {
	byteData, err := json.MarshalIndent(jsonData, "", "    ")
	if err != nil {
		return nil, err
	}

	return byteData, err
}

func ParsingQueryAccount(res string) string {
	strList := strings.Split(res, "sequence")
	strList = strings.Split(strList[1], "\"")
	sequence := strList[2]

	return sequence
}
