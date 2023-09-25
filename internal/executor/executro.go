package executor

import (
	"fmt"
	"os/exec"
)

func Exec(command string, args []string) (string, error) {
	cmd := exec.Command(command, args...)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to exec %s: %w", cmd.String(), err)
	}
	return string(stdout), nil
}
