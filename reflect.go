//+build ignore

package json

import (
	"encoding/json"
)

// Response - модель для тестов
type Response struct {
	Secret string `json:"secret"`
}

// parseJSONByReflect проверяет данные на соответствие json-формату.
func parseJSONByReflect(data []byte) (string, error) {

	var response Response

	err := json.Unmarshal(data, &response)
	if err != nil {
		return "", err
	}

	return response.Secret, nil
}
