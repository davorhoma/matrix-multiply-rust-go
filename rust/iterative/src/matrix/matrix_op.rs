use crate::matrix::{Matrix, matrix::MatrixView};

pub fn sub_vec(a: &MatrixView, b: &MatrixView) -> Matrix {
    let mut result = Matrix::new(a.rows, a.cols);

    for i in 0..a.rows {
        for j in 0..a.cols {
            let idx_a = i * a.stride + j;
            let idx_b = i * b.stride + j;
            result.set(i, j, a.data[idx_a] - b.data[idx_b]);
        }
    }

    result
}