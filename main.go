package main

import (
	"Smarket/router"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", router.Router)

	log.Fatal(http.ListenAndServe(":5500", nil))
}
