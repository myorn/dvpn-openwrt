package node

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

func processByName(name string) (node Node, err error) {
	var out bytes.Buffer

	cmd := exec.Command("ps", "-C", name)
	cmd.Stdout = &out

	err = cmd.Run()

	if err != nil {
		// Process not found
	}

	pidRegExp, err := regexp.Compile("\\w*[0-9]\\w\\w")

	if err != nil {
		return node, err
	}

	pid := pidRegExp.FindString(out.String())

	pidInt, err := strconv.Atoi(pid)

	if err != nil {
		return node, nil
	}

	node.Pid = pidInt

	proc, err := os.FindProcess(pidInt)

	if err != nil {
		return node, err
	}

	timeout := 700 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	node.Online = waitForStateWithContext(ctx, proc)
	return node, err
}

func killProcessByPid(pid int) error {
	proc, err := os.FindProcess(pid)

	if err != nil {
		return err
	}

	return proc.Kill()
}

func waitForStateWithContext(ctx context.Context, proc *os.Process) bool {
	c := make(chan bool)

	go func() {
		state, err := proc.Wait()
		if err == nil {
			c <- state.Exited()
		}

		defer close(c)
	}()

	select {
	case <-ctx.Done():
		// context timed out the process has not exited
		return true
	case state := <-c:
		return state
	}
}
