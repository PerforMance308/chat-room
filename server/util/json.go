package util

import (
	"io/ioutil"
	"encoding/json"
)

func GetJSONConfig(filePath *string, result interface{}) interface{} {
	data, err := ioutil.ReadFile(*filePath)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data, result); err != nil {
		panic(err)
	}

	return result
}
