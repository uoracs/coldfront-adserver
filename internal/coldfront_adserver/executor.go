package coldfront_adserver

import (
	"fmt"
	"io"
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

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("failed to connect to stdout: %v", err)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("failed to connect to stdin: %v", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("failed to connect to stderr: %v", err)
	}

	// go func() {
		// defer stdin.Close()
		fmt.Fprintln(stdin, "Import-Module C:\\racs\\hpcadmin-powershell\\HPCAdmin.psm1 -force")
		fmt.Fprintln(stdin, command)
	// }()

	err = cmd.Start()
	if err != nil {
		return "", fmt.Errorf("failed to start command: %v", err)
	}

	stdoutBytes, err := io.ReadAll(stdoutPipe)
	if err != nil {
		return "", fmt.Errorf("failed to read bytes from stdout: %v", err)
	}
	stdout := DeWindows(string(stdoutBytes))
	stderrBytes, err := io.ReadAll(stderrPipe)
	if err != nil {
		return "", fmt.Errorf("failed to read bytes from stderr: %v", err)
	}
	stderr := DeWindows(string(stderrBytes))
	// all our errors from hpcadmin-powershell are single lines, so lets just
	// grab the first line of stdout
	stderr = strings.Split(stderr, "\n")[0]

	err = cmd.Wait()
	if err != nil {
		return "", fmt.Errorf("command failed: %s. error: %v", stderr, err)
	}

	slog.Debug("command executor", "stdout", stdout)

	// each command includes two output lines from running the stdin commands
	// and one output line at the end for an empty prompt
	// lets just extract the slice between those
	outLines := strings.Split(stdout, "\n")

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
