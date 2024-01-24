package main

import (
	"fmt"
	"net/http"

	"alwadi_markets/router"
)

func main() {
	fs := http.FileServer(http.Dir("./images"))
	http.Handle("/images/", http.StripPrefix("/images/", fs))
	http.HandleFunc("/", router.Router)


	fmt.Println("Starting server at http://192.168.1.5:5500/")
	err := http.ListenAndServe(":5500", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
