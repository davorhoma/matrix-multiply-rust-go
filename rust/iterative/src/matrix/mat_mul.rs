use crate::matrix::{
    matrix::{MatrixView},
    Matrix,
};

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
