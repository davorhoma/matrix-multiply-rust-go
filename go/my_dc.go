package main

import (
	"fmt"
	"sync"
)

func multiplyMatrix(A, B, C Matrix) {
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

func divideAndConquer(A, B, C Matrix) {
	n := len(A)
	if n <= THRESHOLD {
		multiplyMatrix(A, B, C)
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
