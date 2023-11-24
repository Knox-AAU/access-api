package api

import (
	"access-api/pkg/middlewares"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

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
		log.Println("Error opening file:", err)
		return
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Println("Error reading file:", err)
		return
	}

	err = json.Unmarshal(data, &a.services)
	if err != nil {
		log.Println("Error unmarshaling JSON:", err)
		return
	}

	for _, service := range a.services {
		log.Printf("%s: %s\n", service.Name, service.Base_url)
	}
}

func (appState *AppState) LoadEnvs(path string) {
	godotenv.Load(path + ".env")
	appState.internal_key = mustGetEnv("INTERNAL_KEY")
}

func (appState *AppState) ServeHTTP(originRes http.ResponseWriter, originReq *http.Request) {
	url, err := appState.RebuildUrl(originReq.URL.Path)
	if err != nil {
		originRes.Write([]byte("invalid service, got error: " + err.Error() + "\n"))
		return
	}

	var (
		proxyReq *http.Request
	)

	if originReq.Body != nil {
		body, err := io.ReadAll(originReq.Body)
		if err != nil {
			originRes.Write([]byte("Unable to read body from request, got error: " + err.Error() + "\n"))
			return
		}
		proxyReq, err = http.NewRequest(originReq.Method, url, bytes.NewBuffer(body))
	} else {
		proxyReq, err = http.NewRequest(originReq.Method, url, nil)
	}

	if err != nil {
		originRes.Write([]byte("Unable to create request, got error: " + err.Error() + "\n"))
		return
	}

	for key, values := range originReq.Header {
		for _, value := range values {
			proxyReq.Header.Set(key, value)
			log.Printf("Setting header %s as %s\n", key, value)
		}
	}

	log.Printf("Sending proxy request to %s", url)
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

func SetupRouter(path string) http.Handler {
	appState := AppState{
		services: []Service{},
	}
	appState.LoadDataFromFile(path + "services.json")
	appState.LoadEnvs(path)

	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/").HandlerFunc(appState.ServeHTTP)

	router.Use(middlewares.LoggingMiddleware)

	return appState.AuthMiddleware(router)
}
