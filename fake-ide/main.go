package main

import (
	"common/log"
	"fmt"
	"net/http"
	"os"
)

// loggingMiddleware wraps any handler to log requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithField("method", r.Method).
			WithField("path", r.URL.Path).
			WithField("remote", r.RemoteAddr).
			Info("Incoming request")
		next.ServeHTTP(w, r)
	})
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello, World!</h1>")
	fmt.Fprintf(w, "<h2>Environment Variables:</h2><ul>")
	for _, env := range os.Environ() {
		fmt.Fprintf(w, "<li>%s</li>", env)
	}
	fmt.Fprintf(w, "</ul>")
	// Optional: also log to your log system
	for _, env := range os.Environ() {
		log.Info(env)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	// Wrap mux with logging middleware
	loggedMux := loggingMiddleware(mux)
	log.WithField("port", 3000).Info("Starting fake server")
	err := http.ListenAndServe(":3000", loggedMux)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
