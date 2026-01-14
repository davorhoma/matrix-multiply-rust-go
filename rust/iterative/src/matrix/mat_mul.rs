use crate::matrix::{
    matrix::{MatrixView, MatrixViewMut},
    Matrix,
};

pub fn matmul(a: MatrixView, b: MatrixView, c: &mut MatrixViewMut) {
    assert!(a.cols == b.rows, "Incompatible dimensions");

    for i in 0..a.rows {
        for k in 0..a.cols {
            for j in 0..b.cols {
                let v = c.get(i, j) + a.get(i, k) * b.get(k, j);
                c.set(i, j, v);
            }
        }
    }
}

pub fn multiply_matrix(a: MatrixView, b: MatrixView) -> Matrix {
    assert!(a.cols == b.rows, "Incompatible dimensions");

    let mut c = Matrix::new(a.rows, b.cols);

    for i in 0..a.rows {
        for k in 0..a.cols {
            for j in 0..b.cols {
                let v = c.get(i, j) + a.get(i, k) * b.get(k, j);
                c.set(i, j, v);
            }
        }
    }

    c
}

pub fn matrix_matmul(a: Matrix, b: Matrix) -> Matrix {
    assert!(a.cols == b.rows, "Incompatible dimensions");

    let mut c = Matrix::new(a.rows, b.cols);

    for i in 0..a.rows {
        for k in 0..a.cols {
            for j in 0..b.cols {
                let v = c.get(i, j) + a.get(i, k) * b.get(k, j);
                c.set(i, j, v);
            }
        }
    }

    c
}
