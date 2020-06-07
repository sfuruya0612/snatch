package util

import (
	"encoding/base64"
	"fmt"
)

func DecodeString(text string) (string, error) {
	d, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", fmt.Errorf("DecodeString: %v", err)
	}

	return string(d), nil
}
