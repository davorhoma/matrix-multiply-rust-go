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

func multiply_matrix_old(A, B [][]int) [][]int {
	rowsA, colsA := len(A), len(A[0])
	rowsB, colsB := len(B), len(B[0])

	for colsA != rowsB {
		fmt.Println("Matrix multiplication is not possible. Matrices need to be of dimensions: _xM and Mx_")
		return nil
	}

	C := make([][]int, rowsA)
	for i := range C {
		C[i] = make([]int, colsB)
	}

	for i := range rowsA {
		for j := range colsB {
			for k := range colsA {
				C[i][j] += A[i][k] * B[k][j]
			}
		}
	}

	return C
}

func newMatrix(rows, cols int) Matrix {
	C := make(Matrix, rows)
	for i := range C {
		C[i] = make([]int, cols)
	}

	return C
}

func multiplyMatrix(A, B Matrix) Matrix {
	rowsA, colsA := len(A), len(A[0])
	rowsB, colsB := len(B), len(B[0])

	if colsA != rowsB {
		fmt.Println("Incompatible dimensions")
		return nil
	}

	C := newMatrix(rowsA, colsB)
	for i := range rowsA {
		for k := range colsA {
			for j := range colsB {
				C[i][j] += A[i][k] * B[k][j]
			}
		}
	}

	return C
}

func mulDC(a, b, c [][]int) {
	n := len(a)

	// base case
	if n == 1 {
		mulBase(a, b, c)
		return
	}

	// split matrices into views
	a11, a12, a21, a22 := splitView(a)
	b11, b12, b21, b22 := splitView(b)
	c11, c12, c21, c22 := splitView(c)

	// temporary matrices for intermediate products
	t1 := makeZeroMatrix(n / 2)
	t2 := makeZeroMatrix(n / 2)

	// C11 = A11*B11 + A12*B21
	mulDC(a11, b11, t1)
	mulDC(a12, b21, t2)
	add(t1, t2, c11)

	// C12 = A11*B12 + A12*B22
	mulDC(a11, b12, t1)
	mulDC(a12, b22, t2)
	add(t1, t2, c12)

	// C21 = A21*B11 + A22*B21
	mulDC(a21, b11, t1)
	mulDC(a22, b21, t2)
	add(t1, t2, c21)

	// C22 = A21*B12 + A22*B22
	mulDC(a21, b12, t1)
	mulDC(a22, b22, t2)
	add(t1, t2, c22)
}

const THRESHOLD = 128
const PARALLEL_THRESHOLD = 128 //256

func (C *Matrix) multiply(A, B [][]int) {
	rowsA, colsA := len(A), len(A[0])
	rowsB, colsB := len(B), len(B[0])

	if colsA != rowsB {
		fmt.Println("Incompatible dimensions")
		return
	}

	for i := range rowsA {
		for k := range colsA {
			for j := range colsB {
				(*C)[i][j] += A[i][k] * B[k][j]
			}
		}
	}
}

func main() {
	// rows, cols := 2048, 2048

	// A := generate_matrix(rows, cols)
	// B := generate_matrix(rows, cols)

	// C := newMatrix(rows, cols)

	// startSeq := time.Now()
	// D := multiply_matrix_old(A, B)
	// elapsedSeq := time.Since(startSeq)
	// println("Time elapsed seq: ", elapsedSeq.Milliseconds())

	// startSeq2 := time.Now()
	// D2 := multiplyMatrix(A, B)
	// elapsedSeq2 := time.Since(startSeq2)
	// println("Time elapsed multiplyMatrix: ", elapsedSeq2.Milliseconds())

	// if equalMatrices(D, D2) {
	// 	println("Validation passed: matrices are equal")
	// } else {
	// 	println("Validation failed: matrices are NOT equal")
	// }

	// startSeq3 := time.Now()
	// D3 := newMatrix(rows, cols)
	// multiplyMatrixIterative(A, B, D3)
	// elapsedSeq3 := time.Since(startSeq3)
	// println("Time elapsed multiplyMatrixIterative: ", elapsedSeq3.Milliseconds())

	// if equalMatrices(D, D3) {
	// 	println("Validation passed: matrices are equal")
	// } else {
	// 	println("Validation failed: matrices are NOT equal")
	// }

	// startSeq4 := time.Now()
	// D4 := newMatrix(rows, cols)
	// D4.multiply(A, B)
	// elapsedSeq4 := time.Since(startSeq4)
	// println("Time elapsed interface: ", elapsedSeq4.Milliseconds())

	// if equalMatrices(D, D3) {
	// 	println("Validation passed: matrices are equal")
	// } else {
	// 	println("Validation failed: matrices are NOT equal")
	// }

	// startSeq5 := time.Now()
	// D5 := newMatrix(rows, cols)
	// multiplyMatrixIterative_test(A, B, &D5)
	// elapsedSeq5 := time.Since(startSeq5)
	// println("Time elapsed test: ", elapsedSeq5.Milliseconds())

	// if equalMatrices(D, D3) {
	// 	println("Validation passed: matrices are equal")
	// } else {
	// 	println("Validation failed: matrices are NOT equal")
	// }

	// startDC := time.Now()
	// myDCParallel_2(A, B, C)
	// elapsedDC := time.Since(startDC)
	// println("Time elapsed DC: ", elapsedDC.Milliseconds())

	// if equalMatrices(C, D) {
	// 	println("Validation passed: matrices are equal")
	// } else {
	// 	println("Validation failed: matrices are NOT equal")
	// }

	// DC_3 := make([][]int, rows)
	// for i := range C {
	// 	DC_3[i] = make([]int, cols)
	// }

	// startDC_3 := time.Now()
	// divideAndConquer(A, B, DC_3)
	// elapsedDC_3 := time.Since(startDC_3)
	// println("Time elapsed DC_3: ", elapsedDC_3.Milliseconds())

	// if equalMatrices(C, DC_3) {
	// 	println("Validation passed: matrices are equal")
	// } else {
	// 	println("Validation failed: matrices are NOT equal")
	// }

	// S := newMatrix(rows, cols)

	// startStrassen := time.Now()
	// strassen(A, B, S)
	// elapsedStrassen := time.Since(startStrassen)
	// println("Time elapsed Strassen: ", elapsedStrassen.Milliseconds())
	// if equalMatrices(C, S) {
	// 	println("Validation passed: matrices are equal")
	// } else {
	// 	println("Validation failed: matrices are NOT equal")
	// }

	// S2 := newMatrix(rows, cols)

	// startStrassen2 := time.Now()
	// strassen(A, B, S2)
	// elapsedStrassen2 := time.Since(startStrassen2)
	// println("Time elapsed Strassen2: ", elapsedStrassen2.Milliseconds())
	// if equalMatrices(C, S2) {
	// 	println("Validation passed: matrices are equal")
	// } else {
	// 	println("Validation failed: matrices are NOT equal")
	// }

	sizes := []int{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096}
	for _, size := range sizes {
		// fmt.Print(size)
		// A := generate_matrix(size, size)
		// B := generate_matrix(size, size)

		// calculateStrassenTime(size, size, A, B)
		// benchmarkIterative(size)
		// benchmarkDivideAndConquer(size)
		benchmarkStrassen(size, size)
	}
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
		multiplyMatrixIterative(A, B, C)
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
