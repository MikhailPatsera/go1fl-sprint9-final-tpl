package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	SIZE   = 100_000_000
	CHUNKS = 8
)

// generateRandomElements generates a slice of random non-negative ints of the given size.
// Returns nil for non-positive sizes.
func generateRandomElements(size int) []int {
	if size <= 0 {
		return nil
	}

	data := make([]int, size)
	for i := range data {
		data[i] = rand.Int()
	}
	return data
}

// maximum returns the maximum value in the slice.
// Returns 0 if the slice is empty.
func maximum(data []int) int {
	if len(data) == 0 {
		return 0
	}

	max := data[0]
	for _, v := range data[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

// maxChunks splits the slice into CHUNKS parts, finds the maximum in each part
// concurrently, and returns the overall maximum.
func maxChunks(data []int) int {
	if len(data) == 0 {
		return 0
	}
	// Если элементов меньше, чем частей — делим на меньшее число частей.
	if len(data) < CHUNKS {
		return maximum(data)
	}

	chunkSize := len(data) / CHUNKS
	maxes := make([]int, CHUNKS)

	var wg sync.WaitGroup
	wg.Add(CHUNKS)

	for i := 0; i < CHUNKS; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if i == CHUNKS-1 {
			end = len(data) // последняя часть забирает остаток
		}

		go func(idx, from, to int) {
			defer wg.Done()
			maxes[idx] = maximum(data[from:to])
		}(i, start, end)
	}

	wg.Wait()
	return maximum(maxes)
}

func main() {
	fmt.Printf("Генерируем %d целых чисел\n", SIZE)
	data := generateRandomElements(SIZE)

	// Один поток
	fmt.Println("Ищем максимальное значение в один поток")
	start := time.Now()
	max := maximum(data)
	elapsed := time.Since(start).Microseconds()
	fmt.Printf("Максимальное значение элемента: %d\nВремя поиска: %d мкс\n", max, elapsed)

	// CHUNKS потоков
	fmt.Printf("Ищем максимальное значение в %d потоков\n", CHUNKS)
	start = time.Now()
	max = maxChunks(data)
	elapsed = time.Since(start).Microseconds()
	fmt.Printf("Максимальное значение элемента: %d\nВремя поиска: %d мкс\n", max, elapsed)
}