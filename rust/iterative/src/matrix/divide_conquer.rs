use crate::matrix::mat_mul::multiply_matrix;
use crate::matrix::matrix::add_vec;
use crate::matrix::matrix::merge_quadrants;
use crate::matrix::matrix::split_view;
use crate::matrix::matrix::MatrixView;
use crate::matrix::matrix::THRESHOLD;
use crate::matrix::Matrix;

pub fn divide_and_conquer(a: MatrixView, b: MatrixView) -> Matrix {
    let n = a.rows;

    if n <= THRESHOLD {
        return multiply_matrix(a, b);
    }

    let (c11, c12, c21, c22) = compute_blocks_parallel(a, b);
    merge_quadrants(&c11, &c12, &c21, &c22)
}

pub fn dc_sequential(a: MatrixView, b: MatrixView) -> Matrix {
    let n = a.rows;

    if n <= THRESHOLD {
        return multiply_matrix(a, b);
    }

    let (a11, a12, a21, a22) = split_view(a);
    let (b11, b12, b21, b22) = split_view(b);

    use rayon::join;
    let (t1v, t2v) = join(|| divide_and_conquer(a11, b11), || divide_and_conquer(a12, b21));
    let c11 = add_vec(&t1v, &t2v);

    let (t3v, t4v) = join(|| divide_and_conquer(a11, b12), || divide_and_conquer(a12, b22));
    let c12 = add_vec(&t3v, &t4v);

    let (t5v, t6v) = join(|| divide_and_conquer(a21, b11), || divide_and_conquer(a22, b21));
    let c21 = add_vec(&t5v, &t6v);

    let (t7v, t8v) = join(|| divide_and_conquer(a21, b12), || divide_and_conquer(a22, b22));
    let c22 = add_vec(&t7v, &t8v);

    let half = n / 2;
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

fn compute_blocks_parallel(a: MatrixView, b: MatrixView) -> (Matrix, Matrix, Matrix, Matrix) {
    let (a11, a12, a21, a22) = split_view(a);
    let (b11, b12, b21, b22) = split_view(b);

    let ((c11, c12), (c21, c22)) = rayon::join(
        || rayon::join(
            || compute_block(a11, b11, a12, b21),
            || compute_block(a11, b12, a12, b22)),
        || rayon::join(
            || compute_block(a21, b11, a22, b21),
            || compute_block(a21, b12, a22, b22))
    );

    (c11, c12, c21, c22)
}

fn compute_block(a1: MatrixView, b1: MatrixView, a2: MatrixView, b2: MatrixView) -> Matrix {
    let (t1, t2) = rayon::join(
        || divide_and_conquer(a1, b1),
        || divide_and_conquer(a2, b2),
    );
    add_vec(&t1, &t2)
}