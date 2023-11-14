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

	s := systemTest{
		router: router,
		t:      t,
	}
	internalKey := mustGetENV("INTERNAL_KEY")
	header := map[string]string{
		"Access-Authorization": internalKey,
	}
	testName := "TestSystem"
	response := s.sendRequest(http.MethodGet, testName, header)
	require.NotEmpty(t, response, testName)
}

func (s *systemTest) sendRequest(method, testName string, header map[string]string) string {
	const route = "/test"

	req, err := http.NewRequest(string(method), string(route), nil)
	require.NoError(s.t, err, testName, method, route)

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
	return responseData
}
