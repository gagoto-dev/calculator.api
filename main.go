package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
)

var PORT int = 8090

func hello(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"hello\": \"world\"}"))
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

type NumbersRequest struct {
	Numbers []float32 `json:"numbers"`
}

func AddNumbers(numbers []float32) float32 {
	var sum float32

	for _, val := range numbers {
		sum += val
	}
	return sum
}

func handleAddNumbers(w http.ResponseWriter, req *http.Request) {
	var numbers NumbersRequest

	err := json.NewDecoder(req.Body).Decode(&numbers)
	if err != nil {
		w.Write([]byte("Invalid body"))
	}

	sum := AddNumbers(numbers.Numbers)

	w.Write([]byte(fmt.Sprintf("{\"result\": %g}", sum)))
}

func main() {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            true,
	})
	router := http.NewServeMux()

	router.HandleFunc("GET /hello", hello)
	router.HandleFunc("GET /headers", headers)
	router.HandleFunc("POST /addNumbers", handleAddNumbers)

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", PORT),
		Handler: c.Handler(router),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("HTTP Server listening on port %d...\n", PORT)
}
