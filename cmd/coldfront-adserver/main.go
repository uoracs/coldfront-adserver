package main

import (
	"context"
	"log"
	"net/http"
	"os"

	cf "github.com/uoracs/coldfront-adserver/internal/coldfront_adserver"
)

var version string

func main() {
	log.Printf("Starting Coldfront ADServer Version: %s", version)

	executor := cf.NewPowerShellExecutor()

	debug := false
	_, found := os.LookupEnv("DEBUG")
	if found {
		debug = true
	}
	if debug {
		log.Println("Debug mode")
		executor = cf.NewDebugExecutor()
	}

	apiKey, found := os.LookupEnv("COLDFRONT_ADSERVER_API_KEY")
	if !found {
		log.Fatal("You must set COLDFRONT_ADSERVER_API_KEY")
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, cf.DebugKey, debug)
	ctx = context.WithValue(ctx, cf.ApiKeyKey, apiKey)
	ctx = context.WithValue(ctx, cf.ExecutorKey, executor)

	mux := http.NewServeMux()
	mux.HandleFunc("GET 	/", cf.GetOkHandler)
	mux.Handle("POST	/projects", cf.InjectContext(ctx, cf.RequireAuth(http.HandlerFunc(cf.PostProjectsHandler))))
	mux.Handle("GET 	/pstest", cf.InjectContext(ctx, cf.RequireAuth(http.HandlerFunc(cf.GetPSTestHandler))))

	addr := "0.0.0.0:8080"
	log.Printf("Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
