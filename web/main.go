package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fs2 := http.FileServer(http.Dir("./Grid"))
	http.Handle("/Grid/", http.StripPrefix("/Grid/", fs2))
	// port := os.Getenv("port")

	// if port == "" {
	// 	port = "9002"
	// }
	// log.Fatal(http.ListenAndServe(":9002", myRouter))

	log.Println("Listening on :9000...")
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
