package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", corsMiddleware(http.StripPrefix("/static/", fs)))

	fs2 := http.FileServer(http.Dir("./Grid"))
	http.Handle("/Grid/", corsMiddleware(http.StripPrefix("/Grid/", fs2)))

	port := os.Getenv("port")
	if port == "" {
		port = "9000"
	}
	log.Printf("Listening on :%s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// corsMiddleware is a middleware that adds CORS headers to the response
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "GET")
		w.Header().Add("Access-Control-Allow-Headers", "*")

		next.ServeHTTP(w, r)
	})
}
