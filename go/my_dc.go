package main

import (
	"fmt"
	"sync"
)

// func multiplyMatrixIterative_old(A, B, C [][]int) {
// 	n := len(A)
// 	for i := 0; i < n; i++ {
// 		for j := 0; j < n; j++ {
// 			sum := 0
// 			for k := 0; k < n; k++ {
// 				sum += A[i][k] * B[k][j]
// 			}
// 			C[i][j] = sum
// 		}
// 	}
// }

func multiplyMatrixIterative(A, B, C Matrix) {
	rowsA, colsA := len(A), len(A[0])
	rowsB, colsB := len(B), len(B[0])

	if colsA != rowsB {
		fmt.Println("Incompatible dimensions")
		return
	}

	for i := 0; i < rowsA; i++ {
		for k := 0; k < colsA; k++ {
			for j := 0; j < colsB; j++ {
				C[i][j] += A[i][k] * B[k][j]
			}
		}
	}
}

func multiplyMatrixIterative_test(A, B [][]int, C *Matrix) {
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

func myDC(A, B, C [][]int) {
	n := len(A)
	if n <= THRESHOLD {
		multiplyMatrixIterative(A, B, C)
		return
	}

	// split matrices
	a11, a12, a21, a22 := splitView(A)
	b11, b12, b21, b22 := splitView(B)
	c11, c12, c21, c22 := splitView(C)

	// temporary matrices
	t1 := makeZeroMatrix(n / 2)
	t2 := makeZeroMatrix(n / 2)

	// C11 = A11*B11 + A12*B21
	myDC(a11, b11, t1)
	myDC(a12, b21, t2)
	add(t1, t2, c11)

	// C12 = A11*B12 + A12*B22
	myDC(a11, b12, t1)
	myDC(a12, b22, t2)
	add(t1, t2, c12)

	// C21 = A21*B11 + A22*B21
	myDC(a21, b11, t1)
	myDC(a22, b21, t2)
	add(t1, t2, c21)

	// C22 = A21*B12 + A22*B22
	myDC(a21, b12, t1)
	myDC(a22, b22, t2)
	add(t1, t2, c22)
}

func myDCParallel(A, B, C [][]int) {
	n := len(A)
	if n == 1 {
		C[0][0] = A[0][0] * B[0][0]
		return
	}

	// split matrices
	// start := time.Now()
	a11, a12, a21, a22 := splitView(A)
	b11, b12, b21, b22 := splitView(B)
	c11, c12, c21, c22 := splitView(C)
	// fmt.Println("Time elapsed splitView: ", time.Since(start).Milliseconds())

	// temporary matrices
	// start = time.Now()
	t1 := makeZeroMatrix(n / 2)
	t2 := makeZeroMatrix(n / 2)
	// fmt.Println("Time elapsed 2 makeZeroMatrix: ", time.Since(start).Milliseconds())

	var wg sync.WaitGroup
	wg.Add(2)

	// start = time.Now()
	go func() { myDC(a11, b11, t1); wg.Done() }()
	go func() { myDC(a12, b21, t2); wg.Done() }()
	wg.Wait()
	// fmt.Println("TIme elapsed 2 go func(): ", time.Since(start).Milliseconds())
	// start = time.Now()
	add(t1, t2, c11)

	// C12 = A11*B12 + A12*B22
	t1 = makeZeroMatrix(n / 2)
	t2 = makeZeroMatrix(n / 2)
	wg.Add(2)
	go func() { myDC(a11, b12, t1); wg.Done() }()
	go func() { myDC(a12, b22, t2); wg.Done() }()
	wg.Wait()
	add(t1, t2, c12)

	// C21 = A21*B11 + A22*B21
	t1 = makeZeroMatrix(n / 2)
	t2 = makeZeroMatrix(n / 2)
	wg.Add(2)
	go func() { myDC(a21, b11, t1); wg.Done() }()
	go func() { myDC(a22, b21, t2); wg.Done() }()
	wg.Wait()
	add(t1, t2, c21)

	// C22 = A21*B12 + A22*B22
	t1 = makeZeroMatrix(n / 2)
	t2 = makeZeroMatrix(n / 2)
	wg.Add(2)
	go func() { myDC(a21, b12, t1); wg.Done() }()
	go func() { myDC(a22, b22, t2); wg.Done() }()
	wg.Wait()
	add(t1, t2, c22)
}

func myDCParallel_2(A, B, C [][]int) {
	n := len(A)
	if n <= THRESHOLD {
		multiplyMatrixIterative(A, B, C)
		return
	}

	// split matrices
	// start := time.Now()
	a11, a12, a21, a22 := splitView(A)
	b11, b12, b21, b22 := splitView(B)
	c11, c12, c21, c22 := splitView(C)
	// fmt.Println("Time elapsed splitView: ", time.Since(start).Milliseconds())

	// temporary matrices
	// start = time.Now()
	t1 := makeZeroMatrix(n / 2)
	t2 := makeZeroMatrix(n / 2)
	// fmt.Println("Time elapsed 2 makeZeroMatrix: ", time.Since(start).Milliseconds())

	var wg sync.WaitGroup
	wg.Add(2)

	// start = time.Now()
	go func() { myDCParallel_2(a11, b11, t1); wg.Done() }()
	go func() { myDCParallel_2(a12, b21, t2); wg.Done() }()
	wg.Wait()
	// fmt.Println("TIme elapsed 2 go func(): ", time.Since(start).Milliseconds())
	// start = time.Now()
	add(t1, t2, c11)

	// C12 = A11*B12 + A12*B22
	t1 = makeZeroMatrix(n / 2)
	t2 = makeZeroMatrix(n / 2)
	wg.Add(2)
	go func() { myDCParallel_2(a11, b12, t1); wg.Done() }()
	go func() { myDCParallel_2(a12, b22, t2); wg.Done() }()
	wg.Wait()
	add(t1, t2, c12)

	// C21 = A21*B11 + A22*B21
	t1 = makeZeroMatrix(n / 2)
	t2 = makeZeroMatrix(n / 2)
	wg.Add(2)
	go func() { myDCParallel_2(a21, b11, t1); wg.Done() }()
	go func() { myDCParallel_2(a22, b21, t2); wg.Done() }()
	wg.Wait()
	add(t1, t2, c21)

	// C22 = A21*B12 + A22*B22
	t1 = makeZeroMatrix(n / 2)
	t2 = makeZeroMatrix(n / 2)
	wg.Add(2)
	go func() { myDCParallel_2(a21, b12, t1); wg.Done() }()
	go func() { myDCParallel_2(a22, b22, t2); wg.Done() }()
	wg.Wait()
	add(t1, t2, c22)
}

func divideAndConquer(A, B, C Matrix) {
	n := len(A)
	if n <= THRESHOLD {
		multiplyMatrixIterative(A, B, C)
	} else {
		computeBlocksParallel(A, B, C, n)
	}
}

func computeBlocksParallel(A, B, C Matrix, n int) {
	a11, a12, a21, a22 := splitView(A)
	b11, b12, b21, b22 := splitView(B)
	c11, c12, c21, c22 := splitView(C)

	var wg sync.WaitGroup
	wg.Add(4)
	go func() { computeBlock(a11, b11, a12, b21, c11, n/2); wg.Done() }()
	go func() { computeBlock(a11, b12, a12, b22, c12, n/2); wg.Done() }()
	go func() { computeBlock(a21, b11, a22, b21, c21, n/2); wg.Done() }()
	go func() { computeBlock(a21, b12, a22, b22, c22, n/2); wg.Done() }()
	wg.Wait()
}

func computeBlock(a1, b1, a2, b2, cBlock Matrix, size int) {
	t1 := makeZeroMatrix(size)
	t2 := makeZeroMatrix(size)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() { divideAndConquer(a1, b1, t1); wg.Done() }()
	go func() { divideAndConquer(a2, b2, t2); wg.Done() }()
	wg.Wait()

	add(t1, t2, cBlock)
}
