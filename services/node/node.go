package node

import (
	"io"
	"strings"
	"time"
)

var NodeStdOut io.ReadCloser

var NodeStdErr io.ReadCloser

var StartTime time.Time

func init() {
	outReader := strings.NewReader("out reader")
	NodeStdOut = io.NopCloser(outReader)
	errReader := strings.NewReader("err reader")
	NodeStdErr = io.NopCloser(errReader)
}

type Node struct {
	Online    bool
	Pid       int
	IPAddr    string
	StartTime time.Time
}
