package k8s

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// ExecuteCommand executes the command with args and prints output, errors in standard output, error
func ExecuteCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	fmt.Println("CMD: ", cmd)
	setCommandOutAndError(cmd)
	return cmd.Run()
}

func setCommandOutAndError(cmd *exec.Cmd) {
	var errBuf, outBuf bytes.Buffer
	cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)
	cmd.Stdout = io.MultiWriter(os.Stdout, &outBuf)
}
