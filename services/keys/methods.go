package keys

import (
	"bytes"
	"fmt"
	"github.com/audi70r/dvpn-openwrt/services/node"
	"os/exec"
	"strings"
)

func List() (keys Keys, err error) {
	var out bytes.Buffer

	keys.Keys = make([]Key, 0)

	cmd := exec.Command(node.DVPNNodeExec, node.DVPNNodeKeys, node.DVPNNodeList)
	cmd.Stdout = &out

	err = cmd.Run()

	if err != nil {
		fmt.Println(err.Error())
	}

	keysRings := strings.Split(out.String(), "\n")

	if len(keysRings) < 2 {
		return keys, err
	}

	for _, keyRing := range keysRings {
		keyRingSlice := strings.Fields(keyRing)

		if len(keyRingSlice) == 3 && len(keyRingSlice)%3 == 0 {
			key := Key{
				Name:     keyRingSlice[0],
				Operator: keyRingSlice[1],
				Address:  keyRingSlice[2],
			}

			keys.Keys = append(keys.Keys, key)
		}
	}

	return keys, err
}
