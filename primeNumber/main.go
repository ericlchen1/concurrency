package main

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

const THREADS = 10
const LIMIT int32 = 100000000
const BATCH_SIZE = 100

var currentNum int32 = 3
var numPrimeNumbers int32 = 0

func isPrime(n int32) bool {
	ceil := int32(math.Sqrt(float64(n)))
	for i := int32(3); i <= ceil; i+=2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func main() {
	fmt.Printf("Calculating prime numbers to %d\n", LIMIT)

	wg := sync.WaitGroup{}
	
	startTime := time.Now()
	
	for i := 0; i < THREADS; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				n := atomic.AddInt32(&currentNum, int32(BATCH_SIZE))
				numBatchPrimeNumbers := int32(0)
				for i := n - BATCH_SIZE; i < n && i < LIMIT; i+=2 {
					if isPrime(i) {
						numBatchPrimeNumbers += 1
					}
				}
				atomic.AddInt32(&numPrimeNumbers, numBatchPrimeNumbers)
				if n >= LIMIT {
					break
				}
			}
		}()
	}

	wg.Wait()
	
	if LIMIT >= 2 {
		atomic.AddInt32(&numPrimeNumbers, 1) // Include 2
	}
	
	fmt.Printf("Found %d prime numbers\n", numPrimeNumbers)
	fmt.Printf("Took %s\n", time.Since(startTime))
}
