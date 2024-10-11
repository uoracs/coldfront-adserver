package coldfront_adserver

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

func InjectContext(ctx context.Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqWithContext := r.WithContext(ctx)
		next.ServeHTTP(w, reqWithContext)
	})
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("validating api key")
		goodKey := r.Context().Value(ApiKeyKey).(string)
		providedKey := r.Header.Get("X-API-KEY")
		if providedKey != goodKey {
			slog.Error("invalid api key", "providedKey", providedKey)
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "unauthorized")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func GetProjectsHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("GET /projects")
	fmt.Fprint(w, "get projects")
}

func PostProjectsHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("POST /projects")
	var pr CFProjectsRequest
	err := json.NewDecoder(r.Body).Decode(&pr)
	if err != nil {
		slog.Error("failed to decode projects", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid request body")
	}
	for _, project := range pr.Projects {
		err := ProcessProject(r.Context(), project)
		if err != nil {
			slog.Error("failed to process project", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "failed to process project")
		}
	}
}

func GetPSTestHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("GET /pstest")
	ex := r.Context().Value(ExecutorKey).(Executor)
	command := "write-output 'hi there'"
	output, err := ex.Execute(command)
	if err != nil {
		slog.Error("error running powershell command in http endpoint", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "internal server error")
	}
	fmt.Fprint(w, output)
}

func GetOkHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("GET /")
	fmt.Fprint(w, "ok")
}
