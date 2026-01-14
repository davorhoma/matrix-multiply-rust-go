use crate::matrix::{
    Matrix, mat_mul::multiply_matrix, matrix::{MatrixView, THRESHOLD, add_vec, merge_quadrants, split_view}, matrix_op::sub_vec
};

pub fn strassen_sequential(a: MatrixView, b: MatrixView) -> Matrix {
    let n = a.rows;

    // Base case: small matrix, use standard multiplication
    if n <= THRESHOLD {
        return multiply_matrix(a, b); // your standard iterative matmul
    }

    let half = n / 2;

    // Split A and B into 4 submatrices each
    let a11 = a.sub_view(0, 0, half, half);
    let a12 = a.sub_view(0, half, half, half);
    let a21 = a.sub_view(half, 0, half, half);
    let a22 = a.sub_view(half, half, half, half);

    let b11 = b.sub_view(0, 0, half, half);
    let b12 = b.sub_view(0, half, half, half);
    let b21 = b.sub_view(half, 0, half, half);
    let b22 = b.sub_view(half, half, half, half);

    // Strassen's 7 products
    let m1 = strassen_sequential(
        add_vec(&a11.to_matrix(), &a22.to_matrix()).view(),
        add_vec(&b11.to_matrix(), &b22.to_matrix()).view(),
    );
    let m2 = strassen_sequential(add_vec(&a21.to_matrix(), &a22.to_matrix()).view(), b11);
    let m3 = strassen_sequential(a11, sub_vec(&b12, &b22).view());
    let m4 = strassen_sequential(a22, sub_vec(&b21, &b11).view());
    let m5 = strassen_sequential(add_vec(&a11.to_matrix(), &a12.to_matrix()).view(), b22);
    let m6 = strassen_sequential(
        sub_vec(&a21, &a11).view(),
        add_vec(&b11.to_matrix(), &b12.to_matrix()).view(),
    );
    let m7 = strassen_sequential(
        sub_vec(&a12, &a22).view(),
        add_vec(&b21.to_matrix(), &b22.to_matrix()).view(),
    );

    // Combine into result submatrices
    let c11 = add_vec(&sub_vec(&add_vec(&m1, &m4).view(), &m5.view()), &m7);
    let c12 = add_vec(&m3, &m5);
    let c21 = add_vec(&m2, &m4);
    let c22 = add_vec(&sub_vec(&add_vec(&m1, &m3).view(), &m2.view()), &m6);

    // Merge submatrices into final result
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

pub fn strassen(a: MatrixView, b: MatrixView) -> Matrix {
    let n = a.rows;

    if n <= THRESHOLD {
        return multiply_matrix(a, b);
    }

    let a_blocks = split_view(a);
    let b_blocks = split_view(b);
    let matrices = compute_matrices_parallel(a_blocks, b_blocks);
    combine_result(matrices)
}

fn compute_matrices_parallel(
    a: (MatrixView, MatrixView, MatrixView, MatrixView),
    b: (MatrixView, MatrixView, MatrixView, MatrixView),
) -> [Matrix; 7] {
    let (a11, a12, a21, a22) = a;
    let (b11, b12, b21, b22) = b;

    let mut m1 = Matrix::new(a11.rows, a11.cols);
    let mut m2 = Matrix::new(a11.rows, a11.cols);
    let mut m3 = Matrix::new(a11.rows, a11.cols);
    let mut m4 = Matrix::new(a11.rows, a11.cols);
    let mut m5 = Matrix::new(a11.rows, a11.cols);
    let mut m6 = Matrix::new(a11.rows, a11.cols);
    let mut m7 = Matrix::new(a11.rows, a11.cols);

    rayon::scope(|s| {
        s.spawn(|_| m1 = strassen(
            add_vec(&a11.to_matrix(), &a22.to_matrix()).view(), 
            add_vec(&b11.to_matrix(), &b22.to_matrix()).view()
        ));
        s.spawn(|_| m2 = strassen(
            add_vec(&a21.to_matrix(), &a22.to_matrix()).view(), 
            b11
        ));
        s.spawn(|_| m3 = strassen(
            a11, 
            sub_vec(&b12, &b22).view()
        ));
        s.spawn(|_| m4 = strassen(
            a22, 
            sub_vec(&b21, &b11).view()
        ));
        s.spawn(|_| m5 = strassen(
            add_vec(&a11.to_matrix(), &a12.to_matrix()).view(), 
            b22
        ));
        s.spawn(|_| m6 = strassen(
            sub_vec(&a21, &a11).view(), 
            add_vec(&b11.to_matrix(), &b12.to_matrix()).view()
        ));
        s.spawn(|_| m7 = strassen(
            sub_vec(&a12, &a22).view(), 
            add_vec(&b21.to_matrix(), &b22.to_matrix()).view()
        ));
    });

    [m1, m2, m3, m4, m5, m6, m7]
}

fn combine_result(m: [Matrix; 7]) -> Matrix {
    let c11 = add_vec(&sub_vec(&add_vec(&m[0], &m[3]).view(), &m[4].view()), &m[6]);
    let c12 = add_vec(&m[2], &m[4]);
    let c21 = add_vec(&m[1], &m[3]);
    let c22 = add_vec(&sub_vec(&add_vec(&m[0], &m[2]).view(), &m[1].view()), &m[5]);

    merge_quadrants(&c11, &c12, &c21, &c22)
}