package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"access-api/pkg/middlewares"

	"github.com/gorilla/mux"
)

type Service struct {
	Name     string `json:"name"`
	Base_url string `json:"base_url"`
}

type Services struct {
	Services []Service
}

type AppState struct {
	services     Services
	internal_key string
}

func (s Services) RebuildUrl(url string) (string, error) {
	parts := strings.Split(url, "/")
	if len(parts) == 0 {
		return "", fmt.Errorf("url does not specify service")
	}

	for _, service := range s.Services {
		if service.Name == parts[1] {
			return strings.Replace(url, "/"+parts[1], service.Base_url, 1), nil
		}
	}

	return "", fmt.Errorf("invalid service name")
}

func (s *Services) LoadDataFromFile(path string) {
	// Open the JSON file
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	defer file.Close()

	// Read the JSON data from the file
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	err = json.Unmarshal(data, &s.Services)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}
}

func (appState *AppState) LoadEnvs() {
	appState.internal_key = mustGetEnv("INTERNAL_KEY")
}

func (appState *AppState) ServeHTTP(org_res http.ResponseWriter, org_req *http.Request) {
	url, err := appState.services.RebuildUrl(org_req.URL.Path)
	if err != nil {
		org_res.Write([]byte("Invalid service"))
		return
	}

	// Can this be any smarter? I think it is a lot of lines...
	client := &http.Client{}
	body, err := io.ReadAll(org_req.Body)
	if err != nil {
		// I dont fucking care
	}

	// Creates proxy request
	proxy_req, err := http.NewRequest(org_req.Method, url, bytes.NewBuffer(body))

	// Add headers from org_request to proxy_req
	for key, values := range org_req.Header {
		// Key is the header name, values is a slice of header values
		for _, value := range values {
			proxy_req.Header.Set(key, value)
			fmt.Printf("Setting header %s as %s\n", key, value)
		}
	}

	// Send the request
	// WAAAAY to many lines. Please fix
	fmt.Printf("Sending proxy request to %s", url)
	proxy_res, err := client.Do(proxy_req)
	if err != nil {
		org_res.Write([]byte("Unable to send request"))
		return
	}
	defer proxy_res.Body.Close()

	// call middlewares for f.eks. logging
	middlewares.Middlewares(*proxy_req, *proxy_res, *org_req)

	// Read the response body
	proxy_res_body, err := io.ReadAll(proxy_res.Body)
	if err != nil {
	}

	org_res.Write(proxy_res_body)

}

func (appState AppState) AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if appState.internal_key != r.Header.Get("Access-Authorization") {
			http.Error(w, "Authentication error!", http.StatusForbidden)
			return
		}

		r.Header.Del("Access-Authorization")
		h.ServeHTTP(w, r)
	})
}

func mustGetEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("environment variable %s not set", k)
	}

	return v
}

func main() {
	appState := AppState{
		services: Services{},
	}
	appState.services.LoadDataFromFile("./services.json")
	appState.LoadEnvs()

	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/").HandlerFunc(appState.ServeHTTP)

	fmt.Println("Listening at port 8080..")
	log.Fatal(http.ListenAndServe(":8080", appState.AuthMiddleware(router)))
}
