package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"runtime/pprof"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/d5/go-plotly"
	"github.com/d5/go-plotly/plotly/graph_objs"
	"github.com/dustin/go-humanize"
)

type benchmarkResult struct {
	Name    string
	Time    time.Duration
	Allocs  uint64
	Memory  uint64
}

func runBenchmarks() []benchmarkResult {
	b := testing.B{}
	b.Run("Factorial", BenchmarkFactorial)
	return parseBenchmarkResults(&b)
}

func BenchmarkFactorial(b *testing.B) {
	n := big.NewInt(100)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		calculateFactorial(n)
	}
}

func calculateFactorial(n *big.Int) *big.Int {
	result := big.NewInt(1)
	for i := big.NewInt(2); i.Cmp(n) < 0; i.Add(i, big.NewInt(1)) {
		result.Mul(result, i)
	}
	return result
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

func saveResultsToCSV(results []benchmarkResult, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	if err := writer.Write([]string{"Name", "Time (ns)", "Allocs", "Memory (bytes)"}); err != nil {
		return err
	}

	for _, res := range results {
		timeNS := res.Time.Nanoseconds()
		memory := humanize.Bytes(res.Memory)
		if err := writer.Write([]string{res.Name, strconv.Itoa(int(timeNS)), strconv.Itoa(int(res.Allocs)), memory}); err != nil {
			return err
		}
	}

	return writer.Flush()
}

func createDashboard(results []benchmarkResult) {
	trace := graph_objs.Bar{
		X: []string{},
		Y: []float64{},
	}

	for _, res := range results {
		trace.X = append(trace.X, res.Name)
		trace.Y = append(trace.Y, float64(res.Time.Nanoseconds()))
	}

	data := []graph_objs.Trace{&trace}

	layout := graph_objs.Layout{
		Title:   "Benchmark Results",
		Barmode: "group",
		Xaxis: &graph_objs.Xaxis{
			Title: "Benchmark Name",
		},
		Yaxis: &graph_objs.Yaxis{
			Title: "Execution Time (ns)",
		},
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

func serveDashboard(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "high_precision_benchmark_results.html")
}

func main() {
	outputFile := flag.String("output", "benchmark_results.csv", "Output CSV file path")
	flag.Parse()

	results := runBenchmarks()

	// Sort results by time
	sort.Slice(results, func(i, j int) bool {
		return results[i].Time.Seconds() < results[j].Time.Seconds()
	})

	if err := saveResultsToCSV(results, *outputFile); err != nil {
		fmt.Println("Error saving results to CSV:", err)
		return
	}

	createDashboard(results)

	fmt.Println("Dashboard generated. Serving at http://localhost:8080")
	http.HandleFunc("/", serveDashboard)
	http.ListenAndServe(":8080", nil)
}