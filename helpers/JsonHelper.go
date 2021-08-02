package helpers

import "encoding/json"

func ToJson(data interface{}) string {

	bytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}
