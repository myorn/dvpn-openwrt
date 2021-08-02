package dvpnconf

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/pelletier/go-toml"
)

func GetConfig() (config []byte, err error) {
	var dVPNConfig dVPNConfig

	configPath := os.Getenv("HOME") + dVPNConfigRootDir
	confBytes, readErr := ioutil.ReadFile(configPath)

	if readErr != nil {
		return config, readErr
	}

	if err = toml.Unmarshal(confBytes, &dVPNConfig); err != nil {
		return config, err
	}

	config, _ = json.Marshal(dVPNConfig)

	return config, err
}

func PostConfig(config dVPNConfig) (resp []byte, err error) {
	configPath := os.Getenv("HOME") + dVPNConfigRootDir

	configBytes, err := toml.Marshal(config)

	if err != nil {
		return resp, err
	}

	if err = ioutil.WriteFile(configPath, configBytes, 0644); err != nil {
		return resp, err
	}

	resp, err = json.Marshal(config)

	if err != nil {
		return resp, err
	}

	return resp, err
}
