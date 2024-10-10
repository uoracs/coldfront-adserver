package coldfront_adserver

import (
	"fmt"
	"os/exec"
)

func RunPowerShellCommand(c string) (string, error) {
	cmd := exec.Command("powershell", "-nologo", "-noprofile")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("failed to open powershell: %v", err)
	}
	go func() {
		defer stdin.Close()
		fmt.Fprintln(stdin, c)
	}()
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command failed: %v", err)
	}
	return string(out), nil
}
