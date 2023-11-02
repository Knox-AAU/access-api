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

type AppState struct {
	services     []Service
	internal_key string
}

func (a AppState) RebuildUrl(url string) (string, error) {
	parts := strings.Split(url, "/")
	if len(parts) == 0 {
		return "", fmt.Errorf("url does not specify service")
	}

	for _, service := range a.services {
		if service.Name == parts[1] {
			return strings.Replace(url, "/"+parts[1], service.Base_url, 1), nil
		}
	}

	return "", fmt.Errorf("invalid service name")
}

func (a *AppState) LoadDataFromFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	err = json.Unmarshal(data, &a.services)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}
}

func (appState *AppState) LoadEnvs() {
	appState.internal_key = mustGetEnv("INTERNAL_KEY")
}

func (appState *AppState) ServeHTTP(originRes http.ResponseWriter, originReq *http.Request) {
	url, err := appState.RebuildUrl(originReq.URL.Path)
	if err != nil {
		originRes.Write([]byte("invalid service, got error: " + err.Error() + "\n"))
		return
	}

	body, err := io.ReadAll(originReq.Body)
	if err != nil {
		originRes.Write([]byte("Unable to read body from request, got error: " + err.Error() + "\n"))
		return
	}

	proxyReq, err := http.NewRequest(originReq.Method, url, bytes.NewBuffer(body))
	if err != nil {
		originRes.Write([]byte("Unable to create request, got error: " + err.Error() + "\n"))
		return
	}

	for key, values := range originReq.Header {
		for _, value := range values {
			proxyReq.Header.Set(key, value)
			fmt.Printf("Setting header %s as %s\n", key, value)
		}
	}

	fmt.Printf("Sending proxy request to %s", url)
	client := &http.Client{}

	proxyRes, err := client.Do(proxyReq)
	if err != nil {
		originRes.Write([]byte("Unable to send proxy request, got error: " + err.Error() + "\n"))
		return
	}

	defer proxyRes.Body.Close()

	middlewares.Middlewares(*proxyReq, *proxyRes, *originReq)

	proxy_res_body, err := io.ReadAll(proxyRes.Body)
	if err != nil {
		originRes.Write([]byte("Unable to read body from response, got error: " + err.Error() + "\n"))
		return
	}

	originRes.Write(proxy_res_body)
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
		services: []Service{},
	}
	appState.LoadDataFromFile("./services.json")
	appState.LoadEnvs()

	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/").HandlerFunc(appState.ServeHTTP)

	fmt.Println("Listening at port 8080..")
	log.Fatal(http.ListenAndServe(":8080", appState.AuthMiddleware(router)))
}
