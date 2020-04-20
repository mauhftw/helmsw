package lib

import (
	"os"
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

// CheckHelmswDir check the existence the helmsw installation dir
func CheckHelmswDir(HELM_VERSIONS string, HELM_BINS string) error {

	dirs := []string{}
	dirs = append(dirs, HELM_VERSIONS, HELM_BINS)

	// Checks if helmsw dir exists
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			mkdir := &BashCmd{
				Cmd:  "mkdir",
				Args: []string{"-p", dir},
			}
			_, err = ExecBashCmd(mkdir)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
