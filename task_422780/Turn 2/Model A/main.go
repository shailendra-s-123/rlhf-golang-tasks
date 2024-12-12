package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "time"
    "github.com/gorilla/mux"

)
type BenchmarkResult struct {
    Name     string   `json:"name"`
    Duration time.Duration `json:"duration"`
}
func BenchmarkFactorial(b *testing.B) {
    n := int64(100)
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        calculateFactorial(n)
    }
}

func calculateFactorial(n int64) int64 {
    if n < 0 {
        return 0
    }
    if n == 0 || n == 1 {
        return 1
    }

    result := int64(1)
    for i := 2; i <= n; i++ {
        result *= i
    }
    return result
}

func main() {
    fmt.Println("Running benchmarks...")
    b := testing.B{}
    BenchmarkFactorial(&b)
    fmt.Println("Benchmarks completed.")

    results := []BenchmarkResult{
        {
            Name:     "Factorial",
            Duration: b.Tests()[0].RunTests[0].Duration,
        },
    }
    
    storeResults(results)
    
    r := mux.NewRouter()
    r.HandleFunc("/", homeHandler)
    r.HandleFunc("/dashboard", dashboardHandler)
    
    log.Fatal(http.ListenAndServe(":8080", r))
}

func storeResults(results []BenchmarkResult) {
    resultsFile := filepath.Join("results", "benchmark_results.json")
    data, err := json.Marshal(results)
    if err != nil {
        log.Fatalf("Error encoding data: %v", err)
    }
    
    os.MkdirAll(filepath.Join("results"), 0777)
    err = os.WriteFile(resultsFile, data, 0644)
    if err != nil {
        log.Fatalf("Error writing data to file: %v", err)
    }
    
    fmt.Println("Results stored in", resultsFile)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/index.html")
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/dashboard.html")
}