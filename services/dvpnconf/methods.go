package dvpnconf

import (
	"encoding/json"
	"github.com/audi70r/dvpn-openwrt/services/node"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/pelletier/go-toml"
)

func GetConfigs() (config []byte, err error) {
	var dVPNConfig dVPNConfig

	configPath := os.Getenv("HOME") + dVPNConfigPath
	confBytes, readErr := ioutil.ReadFile(configPath)

	// if config could not be read, attempt to init config
	if readErr != nil {
		return initConfig()
	}

	wireguardConfigPath := os.Getenv("HOME") + dVPNWireguardPath
	_, readErr = ioutil.ReadFile(wireguardConfigPath)

	if readErr != nil {
		return initWireguardConfig()
	}

	if err = toml.Unmarshal(confBytes, &dVPNConfig); err != nil {
		return config, err
	}

	config, _ = json.Marshal(dVPNConfig)

	return config, err
}

func PostConfig(config dVPNConfig) (resp []byte, err error) {
	configPath := os.Getenv("HOME") + dVPNConfigPath

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

func initConfig() (config []byte, err error) {
	cmd := exec.Command(node.DVPNNodeExec, node.DVPNNodeConfig, node.DVPNNodeInit)

	err = cmd.Run()

	if err != nil {
		return config, err
	}

	return GetConfigs()
}

func initWireguardConfig() (config []byte, err error) {
	cmd := exec.Command(node.DVPNNodeExec, node.DVPNNodeWireguard, node.DVPNNodeConfig, node.DVPNNodeInit)

	err = cmd.Run()

	if err != nil {
		return config, err
	}

	return GetConfigs()
}
