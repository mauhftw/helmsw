package lib

import (
	"os/exec"
)

type BashCmd struct {
	Cmd      string
	Args     []string
	ExecPath string
}

// ExecBashCmd Executes Bash CLI commands
func ExecBashCmd(c *BashCmd) (string, error) {

	// Set command and argument options
	cmd := c.Cmd
	cmdArgs := c.Args

	// Execute command
	cmdRun := exec.Command(cmd, cmdArgs...)
	cmdRun.Dir = c.ExecPath

	// Print stdout & stderr
	out, err := cmdRun.CombinedOutput()
	if err != nil {
		return string(out), err
	}

	return string(out), nil
}
