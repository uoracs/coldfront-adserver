package coldfront_adserver

import (
	"fmt"
	"log/slog"
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
		return "", fmt.Errorf("failed to connect to stdin: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("failed to connect to stderr: %v", err)
	}
	defer stderr.Close()

	go func() {
		defer stdin.Close()
		fmt.Fprintln(stdin, "Import-Module C:\\racs\\hpcadmin-powershell\\HPCAdmin.psm1 -force")
		fmt.Fprintln(stdin, command)
	}()

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("command failed: %s. error: %v", stderr, err)
	}

	slog.Debug("command executor", "stdout", out)

	// each command includes two output lines from running the stdin commands
	// and one output line at the end for an empty prompt
	// lets just extract the slice between those
	outLines := strings.Split(string(out), "\n")

	if len(outLines) < 3 {
		return "", fmt.Errorf("insufficient output lines")
	}
	filteredLines := outLines[2 : len(outLines)-1]
	var cleansedLines []string
	for _, l := range filteredLines {
		cleansedLines = append(cleansedLines, strings.TrimSpace(l))
	}

	return strings.Join(cleansedLines, "\n"), nil
}

type DebugExecutor struct{}

func NewDebugExecutor() Executor {
	return DebugExecutor{}
}

func (e DebugExecutor) Execute(command string) (string, error) {
	return fmt.Sprintf("debug executing: %s", command), nil
}
