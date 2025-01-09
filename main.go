package main

import (
	"fmt"
)

func main() {
	// Channel for sending integers from worker goroutines to main
	c := make(chan int)

	// Start 5 goroutines
	for i := 1; i <= 5; i++ {
		go func(idx int) {
			// Send numbers from 0..5*idx
			for num := 0; num <= 5*idx; num++ {
				c <- num
			}
			// Indicate this goroutine is done
			c <- -1
		}(i)
	}

	totalSum := 0  // Will accumulate the sum of all received numbers
	doneCount := 0 // Tracks how many goroutines have signaled completion

	// Keep reading from the channel until all 5 goroutines have sent -1
	for {
		val := <-c
		if val == -1 {
			// A goroutine just finished
			doneCount++
			if doneCount == 5 {
				break
			}
		} else {
			// Accumulate the sum of valid numbers
			totalSum += val
		}
	}

	// Print the final sum after all goroutines have finished
	fmt.Println("Final sum:", totalSum)
}
