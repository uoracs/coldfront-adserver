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
		if !ValidAPIKey(r.Context(), r) {
			slog.Error("invalid api key")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "unauthorized")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func ValidAPIKey(ctx context.Context, r *http.Request) bool {
	goodKey := ctx.Value(ApiKeyKey).(string)
	key := r.Header.Get("X-API-KEY")
	return key == goodKey
}

func PostProjectsHandler(w http.ResponseWriter, r *http.Request) {
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
			slog.Error("failed to get process project", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "failed to get process project")
		}
	}
}

func GetPSTestHandler(w http.ResponseWriter, r *http.Request) {
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
	fmt.Fprint(w, "ok")
}
