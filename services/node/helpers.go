package node

import (
	"bytes"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"sync"
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

	wg := sync.WaitGroup{}
	wg.Add(1)
	timeout := 1 * time.Second

	if waitForStateWithTimeout(&wg, timeout, proc, &node) {
		// Timed out waiting for wait group
	} else {
		// Wait group finished
	}
	// Free at last

	return node, err
}

func killProcessByPid(pid int) error {
	proc, err := os.FindProcess(pid)

	if err != nil {
		return err
	}

	return proc.Kill()
}

func waitForStateWithTimeout(wg *sync.WaitGroup, timeout time.Duration, proc *os.Process, node *Node) bool {
	c := make(chan *os.ProcessState)

	go func() {
		state, err := proc.Wait()
		if err != nil {
			c <- state
		}

		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		node.Online = false
		return false // completed normally
	case <-time.After(timeout):
		node.Online = true
		return true // timed out
	}
}
