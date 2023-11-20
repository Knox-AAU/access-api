package main

import (
	"access-api/pkg/api"
	"fmt"
	"log"
	"net/http"
)

func main() {
	router := api.SetupRouter("./")
	log.Println("Listening at port 8080..")
	log.Fatal(http.ListenAndServe(":8080", router))
	fmt.Println("hello world!")
}
