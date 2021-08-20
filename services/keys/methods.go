package keys

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/audi70r/dvpn-openwrt/services/node"
	"github.com/audi70r/dvpn-openwrt/services/socket"
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

	fmt.Println()

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

func AddRecover(req AddRecoverRequest) (err error) {
	//var out bytes.Buffer

	// Run the dvpn-node keys command, setting the key name and recovering it from a mnemonic
	cmd := exec.Command(node.DVPNNodeExec, node.DVPNNodeKeys, node.DVPNNodeAdd, req.Name, node.DVPNNodeRecover)
	nodeStdErr, _ := cmd.StderrPipe()

	mnemonicBuf := bytes.Buffer{}
	mnemonicBuf.Write([]byte(req.Mnemonic + "\n\r"))
	cmd.Stdin = &mnemonicBuf

	if err = cmd.Start(); err != nil {
		return err
	}

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// need to send out to a channel
	stdOutErr := make(chan string)

	go node.SendCaptureAndReturn(nodeStdErr, stdOutErr)

	cmd.Wait()

	stdErr := <-stdOutErr

	close(stdOutErr)
	socket.Conn.Send([]byte(stdErr))

	if strings.Contains(stdErr, "Error") {
		return errors.New(stdErr)
	}

	return nil
}
