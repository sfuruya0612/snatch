package util

import (
	"encoding/json"
	"fmt"
)

type Response struct {
	Instance_id string
	Status      string
	Output      []string
}

func JParser(in interface{}) ([]Response, error) {
	var res []Response

	bytes, err := json.Marshal(in)
	if err != nil {
		return res, fmt.Errorf("Json Marshal error: %v", err)
	}

	err = json.Unmarshal(bytes, &res)
	if err != nil {
		return res, fmt.Errorf("Json Unmarshal error: %v", err)
	}

	return res, nil
}

func Marshal(in interface{}) ([]byte, error) {
	bytes, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("Json Marshal error: %v", err)
	}

	return bytes, nil
}
