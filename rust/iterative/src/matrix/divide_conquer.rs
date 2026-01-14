use crate::matrix::mat_mul::multiply_matrix;
use crate::matrix::matmul;
use crate::matrix::matrix::add_mut;
use crate::matrix::matrix::add_vec;
use crate::matrix::matrix::add_views;
use crate::matrix::matrix::merge_quadrants;
use crate::matrix::matrix::split_new;
use crate::matrix::matrix::split_view;
use crate::matrix::matrix::zero;
use crate::matrix::matrix::MatrixView;
use crate::matrix::matrix::MatrixViewMut;
use crate::matrix::matrix::THRESHOLD;
use crate::matrix::Matrix;

pub fn my_dc(a: MatrixView, b: MatrixView, c: &mut MatrixViewMut) {
    let n = a.rows;

    if n <= THRESHOLD {
        matmul(a, b, c);
        return;
    }

    let (a11, a12, a21, a22) = split_view(a);
    let (b11, b12, b21, b22) = split_view(b);
    let (mut c11, mut c12, mut c21, mut c22) = split_new(c);

    let mut t1 = vec![0; (n / 2) * (n / 2)];
    let mut t2 = vec![0; (n / 2) * (n / 2)];

    let mut t1v = MatrixViewMut {
        data: &mut t1,
        rows: n / 2,
        cols: n / 2,
        stride: n / 2,
    };
    let mut t2v = MatrixViewMut {
        data: &mut t2,
        rows: n / 2,
        cols: n / 2,
        stride: n / 2,
    };

    // C11
    my_dc(a11, b11, &mut t1v);
    my_dc(a12, b21, &mut t2v);
    add_mut(&t1v, &t2v, &mut c11);

    // C12
    // zero(&mut t1v);
    // zero(&mut t2v);
    my_dc(a11, b12, &mut t1v);
    my_dc(a12, b22, &mut t2v);
    add_mut(&t1v, &t2v, &mut c12);

    // C21
    // zero(&mut t1v);
    // zero(&mut t2v);
    my_dc(a21, b11, &mut t1v);
    my_dc(a22, b21, &mut t2v);
    add_mut(&t1v, &t2v, &mut c21);

    // C22
    // zero(&mut t1v); zero(&mut t2v);
    my_dc(a21, b12, &mut t1v);
    my_dc(a22, b22, &mut t2v);
    add_mut(&t1v, &t2v, &mut c22);
}

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

pub fn dc_view(a: MatrixView, b: MatrixView, c: &mut MatrixViewMut) {
    let n = a.rows;

    // Base case: small enough, do iterative multiplication
    if n <= THRESHOLD {
        matmul(a, b, c);
        return;
    }

    let half = n / 2;
    let (a11, a12, a21, a22) = split_view(a);
    let (b11, b12, b21, b22) = split_view(b);

    // Allocate temporary buffers for partial results
    let mut t1 = vec![0; half * half];
    let mut t2 = vec![0; half * half];

    let mut t1v = MatrixViewMut {
        data: &mut t1,
        rows: half,
        cols: half,
        stride: half,
    };
    let mut t2v = MatrixViewMut {
        data: &mut t2,
        rows: half,
        cols: half,
        stride: half,
    };

    // -------- C11 = A11*B11 + A12*B21 --------
    {
        let mut c11 = MatrixViewMut {
            data: &mut c.data[..],
            rows: half,
            cols: half,
            stride: c.stride,
        };
        dc_view(a11, b11, &mut t1v);
        dc_view(a12, b21, &mut t2v);
        add_views(&t1v, &t2v, &mut c11);
    }

    // -------- C12 = A11*B12 + A12*B22 --------
    {
        zero(&mut t1v);
        zero(&mut t2v);

        let mut c12 = MatrixViewMut {
            data: &mut c.data[half..],
            rows: half,
            cols: half,
            stride: c.stride,
        };
        dc_view(a11, b12, &mut t1v);
        dc_view(a12, b22, &mut t2v);
        add_views(&t1v, &t2v, &mut c12);
    }

    // -------- C21 = A21*B11 + A22*B21 --------
    {
        zero(&mut t1v);
        zero(&mut t2v);

        let mut c21 = MatrixViewMut {
            data: &mut c.data[half * c.stride..],
            rows: half,
            cols: half,
            stride: c.stride,
        };
        dc_view(a21, b11, &mut t1v);
        dc_view(a22, b21, &mut t2v);
        add_views(&t1v, &t2v, &mut c21);
    }

    // -------- C22 = A21*B12 + A22*B22 --------
    {
        zero(&mut t1v);
        zero(&mut t2v);

        let mut c22 = MatrixViewMut {
            data: &mut c.data[half * c.stride + half..],
            rows: half,
            cols: half,
            stride: c.stride,
        };
        dc_view(a21, b12, &mut t1v);
        dc_view(a22, b22, &mut t2v);
        add_views(&t1v, &t2v, &mut c22);
    }
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