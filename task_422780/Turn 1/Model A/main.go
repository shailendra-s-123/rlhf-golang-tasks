package main
import (
    "flag"
    "fmt"
    "math"
    "math/big"
    "runtime/pprof"
    "testing"
    "time"
)

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
func main() {
    runtime.GC()
    fmt.Println("High Precision Calculation and Benchmarking")

    // Run the benchmark
    b := testing.B{}
    fmt.Println("Running Benchmark...")
    b.Run("Factorial", BenchmarkFactorial)
    fmt.Println("Benchmark completed")
}