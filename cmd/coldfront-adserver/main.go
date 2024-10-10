package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	cf "github.com/uoracs/coldfront-adserver/internal/coldfront_adserver"
)

var version string

func main() {
	log.Printf("Starting Coldfront ADServer Version: %s", version)
	logLevel := slog.LevelInfo
	executor := cf.NewPowerShellExecutor()

	debugLevel, debug := os.LookupEnv("DEBUG_LEVEL")
	if debug {
		log.Println("Debug mode")
		logLevel = slog.LevelDebug

		if debugLevel == "2" {
			executor = cf.NewDebugExecutor()
		}
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})))

	apiKey, found := os.LookupEnv("COLDFRONT_ADSERVER_API_KEY")
	if !found {
		log.Fatal("You must set COLDFRONT_ADSERVER_API_KEY")
	}
	slog.Debug("api key", "value", apiKey)

	ctx := context.Background()
	ctx = context.WithValue(ctx, cf.DebugKey, debug)
	ctx = context.WithValue(ctx, cf.ApiKeyKey, apiKey)
	ctx = context.WithValue(ctx, cf.ExecutorKey, executor)

	mux := http.NewServeMux()
	mux.HandleFunc("GET 	/", cf.GetOkHandler)
	mux.Handle("GET		/projects", cf.InjectContext(ctx, cf.RequireAuth(http.HandlerFunc(cf.GetProjectsHandler))))
	mux.Handle("POST	/projects", cf.InjectContext(ctx, cf.RequireAuth(http.HandlerFunc(cf.PostProjectsHandler))))
	mux.Handle("GET 	/pstest", cf.InjectContext(ctx, cf.RequireAuth(http.HandlerFunc(cf.GetPSTestHandler))))

	addr := "0.0.0.0:8080"
	log.Printf("Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
