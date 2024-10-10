package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
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

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello world")
	})
	mux.HandleFunc("GET /powershell", func(w http.ResponseWriter, r *http.Request) {
		command := "write-output 'hi there'"
		output, err := RunPowerShellCommand(command)
		if err != nil {
			slog.Error("error running powershell command in http endpoint", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "internal server error")
		}
		fmt.Fprint(w, output)
	})
	addr := "0.0.0.0:8080"
	fmt.Printf("Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
