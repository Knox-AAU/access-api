package middlewares

import (
	"log"
	"net/http"
	"time"
)

func Middlewares(proxy_req http.Request, proxy_res http.Response, org_req http.Request) {
	// Implelemt stuff like logging or metrics here, idk tbh
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		requestID := GenerateRequestID()
		log.Printf("[%s] %s %s - Request received", requestID, r.Method, r.URL.Path)
		r.Header.Add("X-Request-ID", requestID)
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s %s - Response sent in %v", requestID, r.Method, r.URL.Path, time.Since(startTime))
	})
}

func GenerateRequestID() string {
	return time.Now().Format("20060102-15:04:05")
}

func Authentication(h http.Handler, expectedKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if expectedKey != r.Header.Get("Access-Authorization") {
			http.Error(w, "Access API Authentication error!", http.StatusForbidden)
			return
		}

		r.Header.Del("Access-Authorization")
		h.ServeHTTP(w, r)
	})
}
