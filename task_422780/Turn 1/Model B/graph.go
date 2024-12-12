package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kniren/graph-draw"
	"github.com/kniren/graph-draw/draw"
	"github.com/kniren/graph-draw/layout"
	"github.com/kniren/graph-draw/plot"
	"io/ioutil"
	"log"
	"net/http"
)

type result struct {
	Iteration int    `json:"iteration"`
	Time      float64 `json:"time"`
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func serveData(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("benchmark_results.csv")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	results := make([]result, 0)
	if err := csv.NewReader(bytes.NewReader(data)).SkipRows(1).DecodeAll(&results); err != nil {
		log.Fatalf("Error decoding CSV: %v", err)
	}

	json.NewEncoder(w).Encode(results)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", serveIndex)
	r.HandleFunc("/data", serveData)

	fmt.Println("Listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}