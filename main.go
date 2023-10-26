package main

import (
	"Smarket/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", router.Router)

	log.Fatal(http.ListenAndServe(":5500", nil))
	fmt.Println(fmt.Printf("your server runing on http://localhost:%d", 5050))
}