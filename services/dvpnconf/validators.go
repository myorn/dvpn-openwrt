package dvpnconf

import "encoding/json"

func ValidateAndUnmarshal(req []byte) (config dVPNConfig, err error) {
	if err := json.Unmarshal(req, &config); err != nil {
		return config, nil
	}

	return config, nil
}
