package node

import (
	"io"
	"os"
	"strings"
	"time"
)

var NodeStdOut io.ReadCloser

var NodeStdErr io.ReadCloser

var StartTime time.Time

type Node struct {
	Online    bool
	Pid       int
	StartTime time.Time
	OSProcess *os.Process `json:"-"`
}

var ND Node

func init() {
	outReader := strings.NewReader("out reader")
	NodeStdOut = io.NopCloser(outReader)
	errReader := strings.NewReader("err reader")
	NodeStdErr = io.NopCloser(errReader)
}
