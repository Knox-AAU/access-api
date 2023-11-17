package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func Middlewares(proxy_req http.Request, proxy_res http.Response, org_req http.Request) {
	// Implelemt stuff like logging or metrics here, idk tbh
	fmt.Println("Debug")

}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Save the start time, and the request id in the request context
		startTime := time.Now()
		requestID := time.Now().Format("20060102-15:04:05")
		// Log the request
		log.Printf("[%s] %s %s - Request received", requestID, r.Method, r.URL.Path)
		// Add the request id to the response header
		r.Header.Add("Request-ID", requestID)
		// Call the next handler in the chain
		next.ServeHTTP(w, r)
		// Log the response
		log.Printf("[%s] %s %s - Response sent in %v", requestID, r.Method, r.URL.Path, time.Since(startTime))
	})
}
