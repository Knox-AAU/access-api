package api_test

import (
	"access-api/pkg/api"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type systemTest struct {
	router http.Handler
	t      *testing.T
}

func mustGetENV(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Missing required environment variable %s", key)
	}

	return value
}

func TestSystem(t *testing.T) {
	router := api.SetupRouter("../../")
	systemTest{
		router: router,
		t:      t,
	}.sendRequest()
}

func (s systemTest) sendRequest() {
	const (
		route    = "/test"
		testName = "TestSystem"
	)

	internalKey := mustGetENV("INTERNAL_KEY")
	header := map[string]string{
		"Access-Authorization": internalKey,
	}
	req, err := http.NewRequest(string(http.MethodGet), string(route), nil)

	require.NoError(s.t, err, testName, route)

	for key, value := range header {
		req.Header.Set(key, value)
	}

	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)

	response := rr.Result()
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	require.NoError(s.t, err, testName)
	require.Equal(s.t, http.StatusOK, response.StatusCode, testName)
	responseData := string(responseBody)
	require.NotEmpty(s.t, responseData, testName)
}
