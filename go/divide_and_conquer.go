package main

import "sync"

func splitView(m Matrix) (a, b, c, d Matrix) {
	rows := len(m)
	cols := len(m[0])

	r2 := rows / 2
	c2 := cols / 2

	// Top-left
	a = make(Matrix, r2)
	for i := 0; i < r2; i++ {
		a[i] = m[i][:c2]
	}

	// Top-right
	b = make(Matrix, r2)
	for i := 0; i < r2; i++ {
		b[i] = m[i][c2:]
	}

	// Bottom-left
	c = make(Matrix, rows-r2)
	for i := r2; i < rows; i++ {
		c[i-r2] = m[i][:c2]
	}

	// Bottom-right
	d = make(Matrix, rows-r2)
	for i := r2; i < rows; i++ {
		d[i-r2] = m[i][c2:]
	}

	return
}

func makeZeroMatrix(n int) [][]int {
	m := make([][]int, n)
	for i := range m {
		m[i] = make([]int, n)
	}
	return m
}

func mulBase(a, b, c [][]int) {
	c[0][0] = a[0][0] * b[0][0]
}

func add(a, b, out [][]int) {
	for i := range a {
		for j := range a[i] {
			out[i][j] = a[i][j] + b[i][j]
		}
	}
}

func add_new(a, b [][]int) [][]int {
	out := make([][]int, len(a))
	for i := range out {
		out[i] = make([]int, len(a[0]))
	}

	for i := range a {
		for j := range a[i] {
			out[i][j] = a[i][j] + b[i][j]
		}
	}

	return out
}

func sub(a, b, out [][]int) {
	for i := range a {
		for j := range a[i] {
			out[i][j] = a[i][j] - b[i][j]
		}
	}
}

func sub_new(a, b [][]int) [][]int {
	out := make([][]int, len(a))
	for i := range out {
		out[i] = make([]int, len(a[0]))
	}

	for i := range a {
		for j := range a[i] {
			out[i][j] = a[i][j] - b[i][j]
		}
	}

	return out
}

func mulDCSeq(a, b, c [][]int) {
	n := len(a)
	if n <= 1 {
		c[0][0] = a[0][0] * b[0][0]
		return
	}

	a11, a12, a21, a22 := splitView(a)
	b11, b12, b21, b22 := splitView(b)
	c11, c12, c21, c22 := splitView(c)

	t1 := makeZeroMatrix(n / 2)
	t2 := makeZeroMatrix(n / 2)

	// C11
	mulDCSeq(a11, b11, t1)
	mulDCSeq(a12, b21, t2)
	add(t1, t2, c11)

	// C12
	mulDCSeq(a11, b12, t1)
	mulDCSeq(a12, b22, t2)
	add(t1, t2, c12)

	// C21
	mulDCSeq(a21, b11, t1)
	mulDCSeq(a22, b21, t2)
	add(t1, t2, c21)

	// C22
	mulDCSeq(a21, b12, t1)
	mulDCSeq(a22, b22, t2)
	add(t1, t2, c22)
}

const cutoff = 64

func mulDCPar(a, b, c [][]int) {
	n := len(a)

	// base case
	if n == 1 {
		c[0][0] = a[0][0] * b[0][0]
		return
	}

	// stop parallel recursion
	if n <= cutoff {
		mulDC(a, b, c) // your sequential version
		return
	}

	a11, a12, a21, a22 := splitView(a)
	b11, b12, b21, b22 := splitView(b)
	c11, c12, c21, c22 := splitView(c)

	t1 := makeZeroMatrix(n / 2)
	t2 := makeZeroMatrix(n / 2)
	t3 := makeZeroMatrix(n / 2)
	t4 := makeZeroMatrix(n / 2)

	var wg sync.WaitGroup
	wg.Add(4) // parallelize four quadrants

	go func() {
		defer wg.Done()
		// C11 = A11*B11 + A12*B21
		mulDCSeq(a11, b11, t1)
		mulDCSeq(a12, b21, t2)
		add(t1, t2, c11)
	}()

	go func() {
		defer wg.Done()
		// C12 = A11*B12 + A12*B22
		mulDCSeq(a11, b12, t1)
		mulDCSeq(a12, b22, t2)
		add(t1, t2, c12)
	}()

	go func() {
		defer wg.Done()
		// C21 = A21*B11 + A22*B21
		mulDCSeq(a21, b11, t3)
		mulDCSeq(a22, b21, t4)
		add(t3, t4, c21)
	}()

	go func() {
		defer wg.Done()
		// C22 = A21*B12 + A22*B22
		mulDCSeq(a21, b12, t3)
		mulDCSeq(a22, b22, t4)
		add(t3, t4, c22)
	}()

	wg.Wait()
}
