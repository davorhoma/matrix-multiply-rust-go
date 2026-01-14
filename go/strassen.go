package main

import "sync"

func strassenSequential(A, B, C [][]int) {
	n := len(A)

	if n <= THRESHOLD {
		multiplyMatrixIterative(A, B, C)
		return
	}

	half := n / 2

	a11, a12, a21, a22 := splitView(A)
	b11, b12, b21, b22 := splitView(B)
	c11, c12, c21, c22 := splitView(C)

	// S matrices (half x half)
	s1 := makeZeroMatrix(half)
	s2 := makeZeroMatrix(half)
	s3 := makeZeroMatrix(half)
	s4 := makeZeroMatrix(half)
	s5 := makeZeroMatrix(half)
	s6 := makeZeroMatrix(half)
	s7 := makeZeroMatrix(half)
	s8 := makeZeroMatrix(half)
	s9 := makeZeroMatrix(half)
	s10 := makeZeroMatrix(half)

	// P matrices
	p1 := makeZeroMatrix(half)
	p2 := makeZeroMatrix(half)
	p3 := makeZeroMatrix(half)
	p4 := makeZeroMatrix(half)
	p5 := makeZeroMatrix(half)
	p6 := makeZeroMatrix(half)
	p7 := makeZeroMatrix(half)

	// Compute S matrices
	sub(b12, b22, s1)
	add(a11, a12, s2)
	add(a21, a22, s3)
	sub(b21, b11, s4)
	add(a11, a22, s5)
	add(b11, b22, s6)
	sub(a12, a22, s7)
	add(b21, b22, s8)
	sub(a11, a21, s9)
	add(b11, b12, s10)

	// Compute P matrices
	strassenSequential(a11, s1, p1)
	strassenSequential(s2, b22, p2)
	strassenSequential(s3, b11, p3)
	strassenSequential(a22, s4, p4)
	strassenSequential(s5, s6, p5)
	strassenSequential(s7, s8, p6)
	strassenSequential(s9, s10, p7)

	// Combine into C
	// C11 = P5 + P4 - P2 + P6
	add(p5, p4, c11)
	sub(c11, p2, c11)
	add(c11, p6, c11)

	// C12 = P1 + P2
	add(p1, p2, c12)

	// C21 = P3 + P4
	add(p3, p4, c21)

	// C22 = P5 + P1 - P3 - P7
	add(p5, p1, c22)
	sub(c22, p3, c22)
	sub(c22, p7, c22)
}

func strassenParallel(A, B, C [][]int) {
	n := len(A)

	if n <= THRESHOLD {
		multiplyMatrixIterative(A, B, C)
		return
	}

	half := n / 2

	a11, a12, a21, a22 := splitView(A)
	b11, b12, b21, b22 := splitView(B)
	c11, c12, c21, c22 := splitView(C)

	// S matrices
	s1 := makeZeroMatrix(half)
	s2 := makeZeroMatrix(half)
	s3 := makeZeroMatrix(half)
	s4 := makeZeroMatrix(half)
	s5 := makeZeroMatrix(half)
	s6 := makeZeroMatrix(half)
	s7 := makeZeroMatrix(half)
	s8 := makeZeroMatrix(half)
	s9 := makeZeroMatrix(half)
	s10 := makeZeroMatrix(half)

	// P matrices
	p1 := makeZeroMatrix(half)
	p2 := makeZeroMatrix(half)
	p3 := makeZeroMatrix(half)
	p4 := makeZeroMatrix(half)
	p5 := makeZeroMatrix(half)
	p6 := makeZeroMatrix(half)
	p7 := makeZeroMatrix(half)

	// Compute S
	sub(b12, b22, s1)
	add(a11, a12, s2)
	add(a21, a22, s3)
	sub(b21, b11, s4)
	add(a11, a22, s5)
	add(b11, b22, s6)
	sub(a12, a22, s7)
	add(b21, b22, s8)
	sub(a11, a21, s9)
	add(b11, b12, s10)

	// Parallelize only if large enough
	if n >= PARALLEL_THRESHOLD {
		var wg sync.WaitGroup
		wg.Add(7)

		go func() { strassenParallel(a11, s1, p1); wg.Done() }()
		go func() { strassenParallel(s2, b22, p2); wg.Done() }()
		go func() { strassenParallel(s3, b11, p3); wg.Done() }()
		go func() { strassenParallel(a22, s4, p4); wg.Done() }()
		go func() { strassenParallel(s5, s6, p5); wg.Done() }()
		go func() { strassenParallel(s7, s8, p6); wg.Done() }()
		go func() { strassenParallel(s9, s10, p7); wg.Done() }()

		wg.Wait()
	} else {
		// Sequential fallback
		strassenSequential(a11, s1, p1)
		strassenSequential(s2, b22, p2)
		strassenSequential(s3, b11, p3)
		strassenSequential(a22, s4, p4)
		strassenSequential(s5, s6, p5)
		strassenSequential(s7, s8, p6)
		strassenSequential(s9, s10, p7)
	}

	// Combine
	add(p5, p4, c11)
	sub(c11, p2, c11)
	add(c11, p6, c11)

	add(p1, p2, c12)
	add(p3, p4, c21)

	add(p5, p1, c22)
	sub(c22, p3, c22)
	sub(c22, p7, c22)
}

type Matrix [][]int

func strassen_old(A, B, C Matrix) {
	n := len(A)

	if n <= THRESHOLD {
		multiplyMatrixIterative(A, B, C)
		return
	}

	computeMatricesParallel_old(A, B, C, n)
}

func computeMatricesParallel_old(A, B, C Matrix, n int) {
	a11, a12, a21, a22 := splitView(A)
	b11, b12, b21, b22 := splitView(B)
	c11, c12, c21, c22 := splitView(C)

	half := n / 2

	// S matrices
	s := make([]Matrix, 10)
	for i := range s {
		s[i] = makeZeroMatrix(half)
	}

	// M matrices
	m := make([]Matrix, 7)
	for i := range m {
		m[i] = makeZeroMatrix(half)
	}

	// Compute S
	sub(b12, b22, s[0])
	add(a11, a12, s[1])
	add(a21, a22, s[2])
	sub(b21, b11, s[3])
	add(a11, a22, s[4])
	add(b11, b22, s[5])
	sub(a12, a22, s[6])
	add(b21, b22, s[7])
	sub(a11, a21, s[8])
	add(b11, b12, s[9])

	var wg sync.WaitGroup
	wg.Add(7)
	go func() { strassen_old(a11, s[0], m[0]); wg.Done() }()
	go func() { strassen_old(s[1], b22, m[1]); wg.Done() }()
	go func() { strassen_old(s[2], b11, m[2]); wg.Done() }()
	go func() { strassen_old(a22, s[3], m[3]); wg.Done() }()
	go func() { strassen_old(s[4], s[5], m[4]); wg.Done() }()
	go func() { strassen_old(s[6], s[7], m[5]); wg.Done() }()
	go func() { strassen_old(s[8], s[9], m[6]); wg.Done() }()
	wg.Wait()

	// Combine
	add(m[4], m[3], c11)
	sub(c11, m[1], c11)
	add(c11, m[5], c11)

	add(m[0], m[1], c12)

	add(m[2], m[3], c21)

	add(m[4], m[0], c22)
	sub(c22, m[2], c22)
	sub(c22, m[6], c22)
}

func strassen(A, B, C Matrix) {
	n := len(A)

	if n <= THRESHOLD {
		multiplyMatrixIterative(A, B, C)
	} else {
		matrices := computeMatricesParallel(A, B, n)
		combineResult(matrices, C)
	}
}

func computeMatricesParallel(A, B Matrix, n int) [7]Matrix {
	a11, a12, a21, a22 := splitView(A)
	b11, b12, b21, b22 := splitView(B)

	half := n / 2

	// S matrices
	S := [10]Matrix{}
	for i := range S {
		S[i] = makeZeroMatrix(half)
	}

	// Compute S
	sub(b12, b22, S[0])
	add(a11, a12, S[1])
	add(a21, a22, S[2])
	sub(b21, b11, S[3])
	add(a11, a22, S[4])
	add(b11, b22, S[5])
	sub(a12, a22, S[6])
	add(b21, b22, S[7])
	sub(a11, a21, S[8])
	add(b11, b12, S[9])

	// M matrices
	M := [7]Matrix{}
	for i := range M {
		M[i] = makeZeroMatrix(half)
	}

	var wg sync.WaitGroup
	wg.Add(7)
	go func() { strassen(a11, S[0], M[0]); wg.Done() }()
	go func() { strassen(S[1], b22, M[1]); wg.Done() }()
	go func() { strassen(S[2], b11, M[2]); wg.Done() }()
	go func() { strassen(a22, S[3], M[3]); wg.Done() }()
	go func() { strassen(S[4], S[5], M[4]); wg.Done() }()
	go func() { strassen(S[6], S[7], M[5]); wg.Done() }()
	go func() { strassen(S[8], S[9], M[6]); wg.Done() }()
	wg.Wait()

	return M
}

func combineResult(m [7]Matrix, C Matrix) {
	c11, c12, c21, c22 := splitView(C)

	add(m[4], m[3], c11)
	sub(c11, m[1], c11)
	add(c11, m[5], c11)

	add(m[0], m[1], c12)

	add(m[2], m[3], c21)

	add(m[4], m[0], c22)
	sub(c22, m[2], c22)
	sub(c22, m[6], c22)
}

// func mergeQuadrants(c11, c12, c21, c22 Matrix) Matrix {
// 	half := len(c11)
// 	n := half * 2

// 	C := make([][]int, n)
// 	for i := range C {
// 		C[i] = make([]int, n)
// 	}

// 	for i := range n {
// 		for j := range n {
// 			C[i][j] = c11[i][j]
// 			C[i][j+half] = c12[i][j]
// 			C[i+half][j] = c11[i][j]
// 			C[i+half][j+half] = c11[i][j]
// 		}
// 	}

// 	return C
// }
