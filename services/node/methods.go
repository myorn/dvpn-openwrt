package node

import (
	"encoding/json"
	"fmt"
	"github.com/audi70r/dvpn-openwrt/services/socket"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
)

func StartNodeStd(w http.ResponseWriter) (resp []byte, err error) {
	//cmd := exec.Command(dVPNNodeBin, dVPNNodeStart)
	cmd := exec.Command(DVPNNodeExec, DVPNNodeStart)
	NodeStdOut, _ = cmd.StdoutPipe()
	NodeStdErr, _ = cmd.StderrPipe()
	err = cmd.Start()

	go sendAndCapture(w, NodeStdOut)
	go sendAndCapture(w, NodeStdErr)

	return []byte{}, err
}

func StartNode() (resp []byte, err error) {
	cmd := exec.Command(DVPNNodeExec, DVPNNodeStart)
	var stdout, stderr []byte
	var errStdout, errStderr error
	NodeStdOut, _ = cmd.StdoutPipe()
	NodeStdErr, _ = cmd.StderrPipe()
	err = cmd.Start()

	// cmd.Wait() should be called only after we finish reading
	// from stdoutIn and stderrIn.
	// wg ensures that we finish
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		stdout, errStdout = copyAndCapture(os.Stdout, NodeStdOut)
		wg.Done()
	}()

	stderr, errStderr = copyAndCapture(os.Stderr, NodeStdErr)

	wg.Wait()

	err = cmd.Wait()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	if errStdout != nil || errStderr != nil {
		log.Fatal("failed to capture stdout or stderr\n")
	}
	outStr, errStr := string(stdout), string(stderr)
	output := fmt.Sprintf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)

	return []byte(output), err
}

func GetNode() (resp []byte, err error) {
	node, err := processByName(DVPNNodeExec)

	if err != nil {
		return resp, err
	}

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

	return nil
}

func sendAndCapture(w io.Writer, r io.Reader) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			fmt.Fprintf(w, "%v\n\n", string(d))
			//_, err := w.Write(d)
			socket.Connection.WriteMessage(1, d)
			if err != nil {
				return out, err
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
}

func copyAndCapture(w io.Writer, r io.Reader) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			_, err := w.Write(d)
			if err != nil {
				return out, err
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
}
