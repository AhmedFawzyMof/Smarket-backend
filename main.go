package main

import (
	"alwadi_markets/router"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", router.Router)

	fmt.Println("Starting server at http://localhost:5500/")
	err := http.ListenAndServe("localhost:5500", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
