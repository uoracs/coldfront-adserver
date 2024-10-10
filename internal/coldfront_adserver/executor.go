package coldfront_adserver

import (
	"fmt"
	"os/exec"
	"strings"
)

type Executor interface {
	Execute(command string) (string, error)
}

type PowerShellExecutor struct{}

func NewPowerShellExecutor() Executor {
	return PowerShellExecutor{}
}

func (ps PowerShellExecutor) Execute(command string) (string, error) {
	cmd := exec.Command("powershell", "-nologo", "-noprofile")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("failed to open powershell: %v", err)
	}
	go func() {
		defer stdin.Close()
		fmt.Fprintln(stdin, "Import-Module C:\\racs\\hpcadmin-powershell\\HPCAdmin.psm1 -force")
		fmt.Fprintln(stdin, command)
	}()
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command failed: %v", err)
	}

	// each command includes two output lines from running the stdin commands
	// and one output line at the end for an empty prompt
	// lets just extract the slice between those
	outLines := strings.Split(string(out), "\r\n")

	if len(outLines) < 3 {
		return "", fmt.Errorf("insufficient output lines")
	}
	filteredLines := outLines[2 : len(outLines)-1]

	return strings.Join(filteredLines, "\n"), nil
}

type DebugExecutor struct{}

func NewDebugExecutor() Executor {
	return DebugExecutor{}
}

func (e DebugExecutor) Execute(command string) (string, error) {
	return fmt.Sprintf("debug executing: %s", command), nil
}
