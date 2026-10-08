package main

import (
	"comp/backend/generated"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	port := getenv("PORT", "8080")
	allowedOrigin := getenv("ALLOWED_ORIGIN", "http://localhost:5173") // the SvelteKit dev server

	api := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: &Resolver{}}))
	api.AddTransport(transport.POST{})
	api.Use(extension.Introspection{}) // lets the Playground read the schema

	mux := http.NewServeMux()
	mux.Handle("/graphql", withCORS(allowedOrigin, api))
	mux.Handle("/", playground.Handler("Counter", "/graphql"))

	// ReadHeaderTimeout stops slow clients from holding connections open forever.
	server := &http.Server{Addr: ":" + port, Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	log.Printf("API on http://localhost:%s/graphql, Playground on http://localhost:%s/", port, port)
	log.Fatal(server.ListenAndServe())
}

// withCORS lets pages from allowedOrigin call the API from the browser.
func withCORS(allowedOrigin string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions { // preflight: the headers above are the whole answer
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
