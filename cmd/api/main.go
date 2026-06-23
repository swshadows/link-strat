package main

import (
	"fmt"
	"go-api/internal/router"
	"net/http"
)

func main() {
	router := router.NewRouter()

	fmt.Println("Serving API on http://localhost:8080")

	http.ListenAndServe(":8080", router)
}
