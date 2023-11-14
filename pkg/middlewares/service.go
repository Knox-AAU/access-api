package middlewares

import (
	"fmt"
	"net/http"
)

func Middlewares(proxy_req http.Request, proxy_res http.Response, org_req http.Request) {
	// Implelemt stuff like logging or metrics here, idk tbh
	fmt.Println("Debug")

}
