package main

import (
	"fmt"
	"link-strat/internal/router"
	"net/http"
)

func main() {
	router := router.NewRouter()

	fmt.Println("Serving API on http://localhost:8080")

	http.ListenAndServe(":8080", router)
}
