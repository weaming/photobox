package utils

import "encoding/json"

func Jsonify(p interface{}) (string, error) {
	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return "", nil
	}
	return string(b), nil
}
