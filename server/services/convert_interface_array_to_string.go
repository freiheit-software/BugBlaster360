package services

import (
	"fmt"
)

func ConvertToString(data []map[string]interface{}) []map[string]string {
	result := make([]map[string]string, len(data))

	for i, row := range data {
		result[i] = make(map[string]string)

		for key, value := range row {
			switch v := value.(type) {
			case []byte:
				result[i][key] = string(v)
			default:
				result[i][key] = fmt.Sprintf("%v", value)
			}
		}
	}

	return result
}
