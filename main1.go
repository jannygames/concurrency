// nuo 0-1000 masyvą paskirstyti i 10 gijų po 0-99, 100-199 ir rasti jų vidurkius.
// B dalis: Padaryti papildomą metodą, kuris skaičiuotų sumų sumą ir išsiųstų į main.

package main

import (
	"fmt"
)

// sumOfSums is the Part B: it receives all chunk sums, sums them up, and returns the result.
func sumOfSums(sums []int) int {
	total := 0
	for _, s := range sums {
		total += s
	}
	return total
}

// chunkSum calculates the sum of a sub-slice [start, end) of nums and sends it to c.
func chunkSum(nums []int, start, end int, c chan int) {
	sum := 0
	for _, v := range nums[start:end] {
		sum += v
	}
	c <- sum
}

func main() {
	// We want numbers from 0..999 (1000 elements total).
	count := 1000
	nums := make([]int, count)
	for i := 0; i < count; i++ {
		nums[i] = i
	}

	// We will split the slice into 10 chunks of 100 elements each.
	chunkSize := 100
	goroutineCount := 10

	// Create a buffered channel to collect sums from each goroutine.
	sumChan := make(chan int, goroutineCount)

	// Spawn 10 goroutines, each calculating the sum of its 100-element chunk.
	for i := 0; i < goroutineCount; i++ {
		start := i * chunkSize
		end := start + chunkSize
		go chunkSum(nums, start, end, sumChan)
	}

	// Receive sums from all goroutines.
	sums := make([]int, goroutineCount)
	for i := 0; i < goroutineCount; i++ {
		sums[i] = <-sumChan
	}

	// Print the sum and average for each chunk.
	for i, s := range sums {
		avg := float64(s) / float64(chunkSize)
		fmt.Printf("Chunk %2d (indexes %3d..%3d): sum = %5d, average = %8.2f\n",
			i, i*chunkSize, i*chunkSize+chunkSize-1, s, avg)
	}

	// Part B: Sum of all sums
	total := sumOfSums(sums)
	fmt.Println("\nSum of all chunk sums:", total)
}
