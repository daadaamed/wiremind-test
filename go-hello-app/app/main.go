package main

import (
    "fmt"
    "log"
    "net/http"
    "log/slog"
    "os"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
    firstname:= r.URL.Query().Get("firstname")
    lastname:= r.URL.Query().Get("lastname")
    if firstname == "" || lastname == "" {
		http.Error(w, "missing required query parameters: firstname, lastname", http.StatusBadRequest)
		return
	}
    fmt.Fprintf(w, "Hello %s %s", firstname, lastname)
}

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    slog.SetDefault(logger)

    mux := http.NewServeMux()
    mux.HandleFunc("GET /hello", helloHandler)
    
    mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
    })    
    
    server := &http.Server{
      Addr: ":8080",
      Handler: mux,
    }    
    log.Println("listening on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
    }    
}
