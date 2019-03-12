package main

import (
	"log"
	"net/http"

	"github.com/tryOne/routes"
)

func main() {
	log.Fatal(http.ListenAndServe(":8000", routes.R))
}
