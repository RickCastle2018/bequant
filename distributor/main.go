package main

import (
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	log.Println("distributor is starting at localhost:8080")

	http.HandleFunc("/price", handler)
	http.ListenAndServe(":8080", nil)
}
