package main

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

const THREADS = 10
const LIMIT int32 = 1000000
const BATCH_SIZE = 10000
const START_NUM = 2

var currentNum int32 = START_NUM - BATCH_SIZE
var numPrimeNumbers int32 = 0

func isPrime(n int32) bool {
	ceil := int32(math.Sqrt(float64(n)))
	for i := int32(2); i < ceil; i++ {
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
				for n < n + BATCH_SIZE {
					if n >= LIMIT {
						break
					}
					if isPrime(n) {
						numBatchPrimeNumbers += 1
					}
					n += 1
				}
				atomic.AddInt32(&numPrimeNumbers, numBatchPrimeNumbers)

				if currentNum >= LIMIT {
					break
				}
			}
		}()
	}

	wg.Wait()
	
	fmt.Printf("Found %d prime numbers\n", numPrimeNumbers)
	fmt.Printf("Took %s\n", time.Since(startTime))
}
