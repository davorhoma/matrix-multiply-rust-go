package main

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

func add(a, b, out [][]int) {
	for i := range a {
		for j := range a[i] {
			out[i][j] = a[i][j] + b[i][j]
		}
	}
}

func sub(a, b, out [][]int) {
	for i := range a {
		for j := range a[i] {
			out[i][j] = a[i][j] - b[i][j]
		}
	}
}
