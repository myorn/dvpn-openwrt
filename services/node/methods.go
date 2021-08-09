package node

import (
	"encoding/json"
	"fmt"
	"github.com/audi70r/dvpn-openwrt/services/socket"
	"io"
	"os/exec"
	"runtime"
	"time"
)

func StartNodeStd() (resp []byte, err error) {
	cmd := exec.Command(DVPNNodeExec, DVPNNodeStart)
	NodeStdOut, _ = cmd.StdoutPipe()
	NodeStdErr, _ = cmd.StderrPipe()
	err = cmd.Start()

	go sendAndCapture(NodeStdOut)
	go sendAndCapture(NodeStdErr)
	go cmd.Wait()

	StartTime = time.Now()

	return []byte{}, err
}

func GetNode() (resp []byte, err error) {
	node, err := processByName(DVPNNodeExec)

	if err != nil {
		return resp, err
	}

	fmt.Printf("nr on routines %v \n", runtime.NumGoroutine()) // TODO: for debugging, remove later

	node.StartTime = StartTime

	resp, err = json.Marshal(node)

	if err != nil {
		return resp, err
	}

	return resp, nil
}

func KillNode() (err error) {
	node, err := processByName(DVPNNodeExec)

	if err != nil {
		return err
	}

	if err = killProcessByPid(node.Pid); err != nil {
		return err
	}

	StartTime = time.Time{}

	return nil
}

func sendAndCapture(r io.Reader) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			socket.Connection.WriteMessage(1, d)
			if err != nil {
				break
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			break
		}
	}

	return
}
