package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	cf "github.com/uoracs/coldfront-adserver/internal/coldfront_adserver"
)

var version string

func main() {
	log.Printf("Starting Coldfront ADServer Version: %s", version)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello world")
	})
	mux.HandleFunc("POST /projects", func(w http.ResponseWriter, r *http.Request) {
		var pr cf.CFProjectsRequest
		err := json.NewDecoder(r.Body).Decode(&pr)
		if err != nil {
			slog.Error("failed to decode projects", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "invalid request body")
		}
		for _, project := range pr.Projects {
			err := cf.ProcessProject(project)
			if err != nil {
				slog.Error("failed to get process project", "error", err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "failed to get process project")
			}
		}
	})
	mux.HandleFunc("GET /pstest", func(w http.ResponseWriter, r *http.Request) {
		command := "write-output 'hi there'"
		output, err := cf.RunPowerShellCommand(command)
		if err != nil {
			slog.Error("error running powershell command in http endpoint", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "internal server error")
		}
		fmt.Fprint(w, output)
	})
	addr := "0.0.0.0:8080"
	log.Printf("Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
