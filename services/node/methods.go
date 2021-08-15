package node

import (
	"encoding/json"
	"github.com/audi70r/dvpn-openwrt/services/socket"
	"io"
	"os/exec"
	"time"
)

func StartNodeStd() (resp []byte, err error) {
	cmd := exec.Command(DVPNNodeExec, DVPNNodeStart)
	NodeStdOut, _ = cmd.StdoutPipe()
	NodeStdErr, _ = cmd.StderrPipe()

	if err = cmd.Start(); err != nil {
		return []byte{}, err
	}

	ND.Online = true
	ND.StartTime = time.Now()
	ND.OSProcess = cmd.Process
	ND.Pid = cmd.Process.Pid

	go SendAndCapture(NodeStdOut)
	go SendAndCapture(NodeStdErr)

	// After process ended, reset node to defaults
	go func() {
		cmd.Wait()
		ND = Node{}
	}()

	return []byte{}, err
}

func GetNode() (resp []byte, err error) {
	if err = updateNodeStatus(DVPNNodeExec); err != nil {
		return resp, err
	}

	resp, err = json.Marshal(ND)

	if err != nil {
		return resp, err
	}

	return resp, nil
}

func KillNode() (err error) {
	if err = updateNodeStatus(DVPNNodeExec); err != nil {
		return err
	}

	// node is already dead
	if ND.OSProcess == nil {
		return nil
	}

	if err = killProcessByPid(ND.Pid); err != nil {
		return err
	}

	ND = Node{}

	return nil
}

func SendAndCapture(r io.Reader) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			socket.Conn.Send(d)
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

func SendCaptureAndReturn(r io.Reader, stdOutErr chan string) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
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

	stdOutErr <- string(out)

	return
}
