package keys

import "encoding/json"

func ValidateAndUnmarshal(req []byte) (keys AddRecoverRequest, err error) {
	if err := json.Unmarshal(req, &keys); err != nil {
		return keys, nil
	}

	return keys, nil
}
