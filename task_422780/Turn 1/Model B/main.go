package main

import (
	"flag"
	"fmt"
	"math/big"
	"os"
	"runtime"
	"time"
)

func main() {
	precision := flag.Int("p", 10, "precision of big.Float")
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	b := new(big.Float).SetFloat64(2.71828) // e
	x := new(big.Float).SetFloat64(3.0)

	results := benchmarkExponent(b, x, *precision)

	// Write the results to a file for visualization
	writeResultsToFile("benchmark_results.csv", results)
}

func benchmarkExponent(base *big.Float, exponent *big.Float, precision int) []float64 {
	times := make([]float64, 0, 10)
	for i := 0; i < 10; i++ {
		var result big.Float
		result.SetPrec(uint(precision))
		start := time.Now()
		result.Pow(&result, base, exponent)
		times = append(times, float64(time.Since(start).Nanoseconds())/1e6)
	}
	return times
}

func writeResultsToFile(filename string, times []float64) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = fmt.Fprintln(file, "iteration,time(ms)")
	if err != nil {
		panic(err)
	}
	for i, time := range times {
		_, err = fmt.Fprintln(file, i, time)
		if err != nil {
			panic(err)
		}
	}
}