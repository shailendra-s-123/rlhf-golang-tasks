package main  
import (
    // Standard libraries
    "fmt"
    "path/filepath"
    "sort"
    "testing"
    
    // External libraries
    "github.com/d5/go-plotly/plotly"
    "github.com/d5/go-plotly/plotly/graph_objs"
)
func main() {
    runBenchmarks()
    createDashboard()
}
func runBenchmarks() {
    b := testing.B{}
    b.Run("Factorial", BenchmarkFactorial)
    // Parse benchmark results
    results := parseBenchmarkResults(&b)
    
    // Order results by time
    sort.Slice(results, func(i, j int) bool {
        return results[i].Time.Seconds() < results[j].Time.Seconds()
    })
    
    for _, res := range results {
        fmt.Printf("Benchmark %s: %v\n", res.Name, res.Time)
    }
}
func parseBenchmarkResults(b *testing.B) []benchmarkResult {
    var results []benchmarkResult
    
    for _, result := range b.Tests() {
        for _, tc := range result.RunTests {
            results = append(results, benchmarkResult{
                Name:    fmt.Sprintf("%s/%s", result.Name, tc.Name),
                Time:    tc.Duration,
                Allocs:  tc.Allocs,
                Memory: tc.Memory,
            })
        }
    }
    
    return results
}
type benchmarkResult struct {
    Name    string
    Time    time.Duration
    Allocs  uint64
    Memory uint64
}
func createDashboard() {
    trace := graph_objs.Bar{
        X: []string{"Factorial"},
        Y: []float64{5230.85745}, // Replace this with the actual values
        // Add more traces for additional benchmarks if needed
        // Name: "Precision Operation",
    }
    
    data := []graph_objs.Trace{&trace}
    
    layout := graph_objs.Layout{
        Title: "Benchmark Results",
        Barmode: "group",
    }
    
    fig := graph_objs.Figure{Data: data, Layout: &layout}
    
    plotURL, err := plotly.Plot(fig, "high_precision_benchmark_results.html", true, plotly.Config{
        Offline:    true,
        AutoOpen: false,
    })
    
    if err != nil {
        panic(err)
    }
    
    fmt.Println("Dashboard created:", plotURL)
}