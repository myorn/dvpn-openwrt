package node

import (
	"bytes"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

func updateNodeStatus(name string) (err error) {
	if ND.OSProcess != nil {
		return nil
	}

	// Search for the process in case node was started elsewhere. ex. through ssh console.
	var out bytes.Buffer

	cmd := exec.Command("ps", "-C", name)
	cmd.Stdout = &out

	err = cmd.Run()

	if err != nil {
		return nil
	}

	pidRegExp, err := regexp.Compile("\\w*[0-9]\\w\\w")

	if err != nil {
		return err
	}

	pid := pidRegExp.FindString(out.String())

	pidInt, err := strconv.Atoi(pid)

	if err != nil {
		return err
	}

	ND.OSProcess, err = os.FindProcess(pidInt)

	// process not found, so what...
	if err != nil {
		return nil
	}

	ND.Pid = pidInt
	ND.Online = true

	go func() {
		ND.OSProcess.Wait()
		ND = Node{}
	}()

	return err
}

func killProcessByPid(pid int) error {
	proc, err := os.FindProcess(pid)

	if err != nil {
		return err
	}

	return proc.Kill()
}
