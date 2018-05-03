package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/AirGateway/agw-proxy/logger"
	"github.com/gorilla/mux"
)

const (
	apiURLEnv = "API_URL"
	apiKeyEnv = "API_KEY"
	portEnv   = "PORT"
	listenEnv = "LISTEN"

	defaultAPIURL = "https://proxy.airgateway.net"
	defaultPort   = "8000"
	defaultListen = "0.0.0.0"
)

var (
	apiURL     *url.URL
	apiKey     string
	port       int
	listen     string
	httpClient = new(http.Client)
)

func init() {
	apiURLStr := os.Getenv(apiURLEnv)
	if apiURLStr == "" {
		apiURLStr = defaultAPIURL
	}

	apiURLParsed, err := url.Parse(apiURLStr)
	if err != nil {
		logger.Fatalf("ERROR: URL passed is not valid: %s!", apiURLStr)
	}
	apiURL = apiURLParsed

	apiKey = os.Getenv(apiKeyEnv)
	if apiKey == "" {
		logger.Fatalf("ERROR: An API Key must be specified in %s env var!", apiKeyEnv)
	}

	portStr := os.Getenv(portEnv)
	if portStr == "" {
		portStr = defaultPort
	}

	portInt, err := strconv.Atoi(portStr)
	if err != nil {
		logger.Fatalf("ERROR: The port number specified is not a number: %s!", portStr)
	}
	port = portInt

	listen = os.Getenv(listenEnv)
	if listen == "" {
		listen = defaultListen
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/{version}/{method}", processHandler).Methods("Post")
	router.HandleFunc("/status", statusHandler).Methods("Get")
	router.HandleFunc("/status/", statusHandler).Methods("Get")

	logger.Info("[HTTP] Started listening at %s:%d", listen, port)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", listen, port),
		Handler: logRequest(router),
	}

	// Run our server in a goroutine so that it doesn't block.
	if err := srv.ListenAndServe(); err != nil {
		logger.Panic(err)
	}
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("[HTTP] %s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

/*
Status check
It just return a string with a HTTP 200 response code to verify that
the server is alive.
*/
func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Status: OK")
}

/*
Process function
It forward the request adding the authorisation header to the upstream API.
*/
func processHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	apiMethod := vars["method"]
	apiVersion := vars["version"]

	path := fmt.Sprintf("%s/%s", apiVersion, apiMethod)
	relative, _ := url.Parse(path)
	fullURL := apiURL.ResolveReference(relative).String()

	buffer, _ := ioutil.ReadAll(r.Body)
	req, err := http.NewRequest(r.Method, fullURL, bytes.NewBuffer(buffer))
	if err != nil {
		logger.Error("Error proxying request: ", err)
	}

	req.Header = r.Header
	req.Header.Add("Authorization", apiKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Error("Error proxying request: ", err)
	}

	for k, v := range resp.Header {
		w.Header().Set(k, v[0])
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
	resp.Body.Close()
}
