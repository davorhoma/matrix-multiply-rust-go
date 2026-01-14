#[derive(Debug)]
pub struct Matrix {
    pub rows: usize,
    pub cols: usize,
    pub data: Vec<i32>,
}

impl Matrix {
    pub fn new(rows: usize, cols: usize) -> Self {
        Self {
            rows,
            cols,
            data: vec![0; rows * cols],
        }
    }

    pub fn get(&self, r: usize, c: usize) -> i32 {
        self.data[r * self.cols + c]
    }

    pub fn set(&mut self, r: usize, c: usize, val: i32) {
        self.data[r * self.cols + c] = val;
    }

    pub fn equals(&self, other: &Matrix) -> bool {
        if self.rows != other.rows || self.cols != other.cols {
            return false;
        }

        self.data == other.data
    }

    pub fn view(&'_ self) -> MatrixView<'_> {
        MatrixView {
            data: &self.data[..],
            rows: self.rows,
            cols: self.cols,
            stride: self.cols,
        }
    }
    
}

#[derive(Clone, Copy)]
pub struct MatrixView<'a> {
    pub data: &'a [i32],
    pub rows: usize,
    pub cols: usize,
    pub stride: usize,
}

impl<'a> MatrixView<'a> {
    #[inline]
    pub fn get(&self, r: usize, c: usize) -> i32 {
        self.data[r * self.stride + c]
    }

    pub fn to_matrix(&self) -> Matrix {
        let mut result = Matrix::new(self.rows, self.cols);

        for i in 0..self.rows {
            for j in 0..self.cols {
                let idx = i * self.stride + j;
                result.set(i, j, self.data[idx]);
            }
        }

        result
    }

    pub fn sub_view(&self, row: usize, col: usize, rows: usize, cols: usize) -> MatrixView<'a> {
        MatrixView {
            data: &self.data[row * self.stride + col..],
            rows,
            cols,
            stride: self.stride,
        }
    }
}

pub fn split_view<'a>(m: MatrixView<'a>) -> (MatrixView<'a>, MatrixView<'a>, MatrixView<'a>, MatrixView<'a>) {
    let r2 = m.rows / 2;
    let c2 = m.cols / 2;

    let a11 = MatrixView {
        data: m.data,
        rows: r2,
        cols: c2,
        stride: m.stride,
    };

    let a12 = MatrixView {
        data: &m.data[c2..],
        rows: r2,
        cols: m.cols - c2,
        stride: m.stride,
    };

    let a21 = MatrixView {
        data: &m.data[r2 * m.stride..],
        rows: m.rows - r2,
        cols: c2,
        stride: m.stride,
    };

    let a22 = MatrixView {
        data: &m.data[r2 * m.stride + c2..],
        rows: m.rows - r2,
        cols: m.cols - c2,
        stride: m.stride,
    };

    (a11, a12, a21, a22)
}

pub fn add_vec(a: &Matrix, b: &Matrix) -> Matrix {
    let mut matrix = Matrix::new(a.rows, b.cols);

    for i in 0..a.rows {
        for j in 0..a.cols {
            matrix.set(i, j, a.get(i, j) + b.get(i, j));
        }
    }

    matrix
}

pub const THRESHOLD: usize = 128;

pub fn merge_quadrants(c11: &Matrix, c12: &Matrix, c21: &Matrix, c22: &Matrix) -> Matrix {
    let half = c11.rows;
    let n = half * 2;

    let mut c = Matrix::new(n, n);

    for i in 0..half {
        for j in 0..half {
            c.set(i, j, c11.get(i, j));
            c.set(i, j + half, c12.get(i, j));
            c.set(i + half, j, c21.get(i, j));
            c.set(i + half, j + half, c22.get(i, j));
        }
    }

    c
}

