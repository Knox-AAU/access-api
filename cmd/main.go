package main

import (
	"access-api/pkg/api"
	"fmt"
	"log"
	"net/http"
)

func main() {
	router := api.SetupRouter("./")
	port := "80"
	log.Printf("Listening at port %s..\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
