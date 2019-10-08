package util

import "encoding/json"

type Response struct {
	Instance_id string
	Status      string
	Output      []string
}

func JParser(in interface{}) ([]Response, error) {
	var res []Response

	bytes, err := json.Marshal(in)
	if err != nil {
		return res, err
	}

	json.Unmarshal(bytes, &res)

	return res, nil
}

func Marshal(in interface{}) ([]byte, error) {
	bytes, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
