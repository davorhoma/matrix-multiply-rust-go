package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func generate_matrix(rows, cols int) [][]int {
	if rows == -1 {
		rows = rand.Intn(20) + 1
	}
	if cols == -1 {
		cols = rand.Intn(20) + 1
	}
	matrix := make([][]int, rows)
	for i := 0; i < rows; i++ {
		matrix[i] = make([]int, cols)
	}

	for i, row := range matrix {
		for j := range row {
			matrix[i][j] = rand.Intn(10)
		}
	}

	return matrix
}

func print_matrix(m [][]int) {
	for i, row := range m {
		for j := range row {
			fmt.Printf("%d ", m[i][j])
		}
		fmt.Println()
	}
	fmt.Println()
}

func newMatrix(rows, cols int) Matrix {
	C := make(Matrix, rows)
	for i := range C {
		C[i] = make([]int, cols)
	}

	return C
}

const THRESHOLD = 128
const PARALLEL_THRESHOLD = 128

func main() {
	rows, cols := 2048, 2048

	A := generate_matrix(rows, cols)
	B := generate_matrix(rows, cols)

	startIterative := time.Now()
	C := newMatrix(rows, cols)
	multiplyMatrix(A, B, C)
	elapsedIterative := time.Since(startIterative)
	println("Time elapsed iterative: ", elapsedIterative.Milliseconds())

	startDC := time.Now()
	DC := newMatrix(rows, cols)
	divideAndConquer(A, B, DC)
	elapsedDC := time.Since(startDC)
	println("Time elapsed DC: ", elapsedDC.Milliseconds())

	if equalMatrices(C, DC) {
		println("Validation passed: matrices are equal")
	} else {
		println("Validation failed: matrices are NOT equal")
	}

	S := newMatrix(rows, cols)

	startStrassen := time.Now()
	strassen(A, B, S)
	elapsedStrassen := time.Since(startStrassen)
	println("Time elapsed Strassen: ", elapsedStrassen.Milliseconds())
	if equalMatrices(C, S) {
		println("Validation passed: matrices are equal")
	} else {
		println("Validation failed: matrices are NOT equal")
	}

	// sizes := []int{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096}
	// for _, size := range sizes {
	// 	// fmt.Print(size)
	// 	A := generate_matrix(size, size)
	// 	B := generate_matrix(size, size)

	// 	calculateStrassenTime(size, size, A, B)
	// 	benchmarkIterative(size)
	// 	benchmarkDivideAndConquer(size)
	// 	benchmarkStrassen(size, size)
	// }
}

func equalMatrices(a, b [][]int) bool {
	if len(a) != len(b) || len(a[0]) != len(b[0]) {
		return false
	}
	for i := range a {
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

func calculateStrassenTime(rows, cols int, A, B Matrix) {
	// var strassen_time time.Duration
	var strassen_time_new time.Duration

	for range 20 {
		// S := newMatrix(rows, cols)
		// startStrassen := time.Now()
		// strassen_old(A, B, S)
		// elapsedStrassen := time.Since(startStrassen)
		// // println("Time elapsed Strassen: ", elapsedStrassen.Milliseconds())
		// strassen_time += elapsedStrassen

		S_new := newMatrix(rows, cols)
		startStrassen_new := time.Now()
		strassen(A, B, S_new)
		elapsedStrassen_new := time.Since(startStrassen_new)
		// println("Time elapsed Strassen_new: ", elapsedStrassen_new.Milliseconds())
		strassen_time_new += elapsedStrassen_new

	}

	// println("Total time elapsed Strassen: ", rows, "x", cols, strassen_time.Milliseconds())
	println("Total time elapsed Strassen_new: ", rows, "x", cols, strassen_time_new.Milliseconds())
	// time.Sleep(3 * time.Second)
}

func benchmarkIterative(size int) {
	var iterativeTime time.Duration

	for range 20 {
		A := generate_matrix(size, size)
		B := generate_matrix(size, size)

		start := time.Now()
		C := newMatrix(size, size)
		multiplyMatrix(A, B, C)
		iterativeTime += time.Since(start)
		// println(i)
	}

	avg := iterativeTime / 20

	println("Total time elapsed Strassen:", size, "x", size, iterativeTime.Milliseconds())
	println("Average Strassen:", size, "x", size, avg.Milliseconds())

	// Write to CSV
	// if err := writeAvgToCSV("iterative_go.csv", size, avg); err != nil {
	// 	panic(err)
	// }
}

func benchmarkDivideAndConquer(size int) {
	var dcTime time.Duration

	for i := range 20 {
		A := generate_matrix(size, size)
		B := generate_matrix(size, size)

		start := time.Now()
		C := newMatrix(size, size)
		divideAndConquer(A, B, C)
		dcTime += time.Since(start)
		println(i)
	}

	avg := dcTime / 20

	println("Total time elapsed DC:", size, "x", size, dcTime.Milliseconds())
	println("Average DC:", size, "x", size, avg.Milliseconds())

	// Write to CSV
	// if err := writeAvgToCSV("divideAndConquer_go.csv", size, avg); err != nil {
	// 	panic(err)
	// }
}

func benchmarkStrassen(rows, cols int) {
	var strassenTime time.Duration

	for range 20 {
		A := generate_matrix(rows, cols)
		B := generate_matrix(rows, cols)

		start := time.Now()
		S := newMatrix(rows, cols)
		strassen(A, B, S)
		strassenTime += time.Since(start)
		// println(i)
	}

	avg := strassenTime / 20

	println("Total time elapsed Strassen:", rows, "x", cols, strassenTime.Milliseconds())
	println("Average Strassen:", rows, "x", cols, avg.Milliseconds())

	// Write to CSV
	// if err := writeAvgToCSV("strassen_results.csv", rows, avg); err != nil {
	// 	panic(err)
	// }
}

func writeAvgToCSV(filename string, rows int, avg time.Duration) error {
	// Open file in append mode, create if not exists
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{
		strconv.Itoa(rows),
		strconv.FormatInt(avg.Microseconds(), 10),
	}

	return writer.Write(record)
}
