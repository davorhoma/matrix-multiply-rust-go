mod matrix;

use std::env;

use rand::Rng;
use chrono::offset::Utc;
use matrix::Matrix;

use crate::matrix::divide_conquer::divide_and_conquer;
use crate::matrix::mat_mul::multiply_matrix;
use crate::matrix::matrix::MatrixView;
use crate::matrix::strassen::strassen;
use crate::matrix::strassen::strassen_sequential;

fn generate_matrix(rows: usize, cols: usize) -> Matrix {
    let mut rng = rand::thread_rng();

    let rows = rows as usize;
    let cols = cols as usize;

    let mut data = Vec::with_capacity(rows * cols);
    for _ in 0..rows * cols {
        data.push(rng.gen_range(0..10));
    }

    Matrix { rows, cols, data }
}

fn is_power_of_two(n: usize) -> bool {
    n != 0 && (n & (n - 1)) == 0
}

fn main() {
    let args: Vec<String> = env::args().collect();

    let rows;
    let cols;
    if args.len() != 2 {
        rows = 1024;
        cols = 1024;
    } else {
        rows = args[1].parse().expect("Rows must be a positive integer");
        cols = args[1].parse().expect("Cols must be a positive integer");
    }

    if !is_power_of_two(rows) || !is_power_of_two(cols) {
        eprintln!(
            "Error: rows and cols must be powers of two (got {}x{})",
            rows, cols
        );
        std::process::exit(1);
    }

    let a = generate_matrix(rows, cols);
    let b = generate_matrix(rows, cols);
    
    let a_view = MatrixView { data: &a.data, rows: a.rows, cols: a.cols, stride: a.cols };
    let b_view = MatrixView { data: &b.data, rows: b.rows, cols: b.cols, stride: b.cols };
    
    // -------- Iterative result --------
    let start_iterative = Utc::now().time();
    let c = multiply_matrix(a_view, b_view);
    let end_iterative = Utc::now().time();
    let diff_iterative = end_iterative - start_iterative;
    println!("Iterative multiplication new_matmul: {} ms", diff_iterative.num_milliseconds());
    //-----------------------------------
    
    // -------- Divide & conquer result --------
    let start_dc = Utc::now().time();
    let c_dc = divide_and_conquer(a_view, b_view);
    let end_dc = Utc::now().time();
    let diff_dc = end_dc - start_dc;
    println!("Divide and conquer: {} ms", diff_dc.num_milliseconds());

    // -------- Compare --------
    if c.equals(&c_dc) {
        println!("Matrices are equal");
    } else {
        println!("Matrices are NOT equal");
    }
    //-----------------------------------

    // // ----------- STRASSEN SEQUENTIAL --------------
    // let start_strassen = Utc::now().time();
    // let c_strassen = strassen_sequential(a_view, b_view);
    // let end_strassen = Utc::now().time();
    // let diff_strassen = end_strassen - start_strassen;
    // println!("Strassen: {} ms", diff_strassen.num_milliseconds());

    //  // -------- Compare view --------
    // if c.equals(&Matrix { rows: rows, cols: cols, data: c_strassen.data }) {
    //     println!("Matrices are equal");
    // } else {
    //     println!("Matrices are NOT equal");
    // }

    // ----------- STRASSEN PARALLEL --------------
    let start_strassen_parallel = Utc::now().time();
    let c_strassen_parallel = strassen(a_view, b_view);
    let end_strassen_parallel = Utc::now().time();
    let diff_strassen_parallel = end_strassen_parallel - start_strassen_parallel;
    println!("Strassen parallel: {} ms", diff_strassen_parallel.num_milliseconds());

     // -------- Compare view --------
    if c.equals(&Matrix { rows: rows, cols: cols, data: c_strassen_parallel.data }) {
        println!("Matrices are equal");
    } else {
        println!("Matrices are NOT equal");
    }
    //-----------------------------------

    // ----------- BENCHMARK --------------
    // let n = [1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048];
    // for i in n {
    //     println!("{i}");
    //     benchmark_strassen_parallel_csv(i, "strassen_parallel.csv", false);
    //     benchmark_iterative_csv(i, "iterative_rust.csv", false);
    //     benchmark_dc_csv(i, "dc_rust.csv", false);
    // }
}

fn benchmark_strassen_parallel_csv(
    rows: usize,
    csv_path: &str,
    save: bool,
) {
    let mut total_us: u128 = 0;

    for _ in 0..20 {
        let a = generate_matrix(rows, rows);
        let b = generate_matrix(rows, rows);
        
        let a_view = MatrixView { data: &a.data, rows: a.rows, cols: a.cols, stride: a.cols };
        let b_view = MatrixView { data: &b.data, rows: b.rows, cols: b.cols, stride: b.cols };

        let start = std::time::Instant::now();
        let _c = strassen(a_view, b_view);
        let elapsed = start.elapsed();

        total_us += elapsed.as_micros();
    }

    let avg_us = total_us / 20;

    println!(
        "Strassen parallel avg: {}x{} -> {} µs",
        rows, rows, avg_us
    );

    if save {
        use std::io::Write;
        let mut file = std::fs::OpenOptions::new()
            .create(true)
            .append(true)
            .open(csv_path)
            .expect("Failed to open CSV");
    
        writeln!(file, "{},{}", rows, avg_us)
            .expect("Failed to write CSV");
    }
}

fn benchmark_iterative_csv(
    rows: usize,
    csv_path: &str,
    save: bool,
) {
    let mut total_us: u128 = 0;

    for _ in 0..20 {
        let a = generate_matrix(rows, rows);
        let b = generate_matrix(rows, rows);
        
        let a_view = MatrixView { data: &a.data, rows: a.rows, cols: a.cols, stride: a.cols };
        let b_view = MatrixView { data: &b.data, rows: b.rows, cols: b.cols, stride: b.cols };

        let start = std::time::Instant::now();
        let _c = multiply_matrix(a_view, b_view);
        let elapsed = start.elapsed();

        total_us += elapsed.as_micros();
    }

    let avg_us = total_us / 20;

    println!(
        "Iterative avg: {}x{} -> {} µs",
        rows, rows, avg_us
    );

    if save {
        use std::io::Write;
        let mut file = std::fs::OpenOptions::new()
            .create(true)
            .append(true)
            .open(csv_path)
            .expect("Failed to open CSV");
    
        writeln!(file, "{},{}", rows, avg_us)
            .expect("Failed to write CSV");
    }
}

fn benchmark_dc_csv(
    rows: usize,
    csv_path: &str,
    save: bool
) {
    let mut total_us: u128 = 0;

    for _ in 0..20 {
        let a = generate_matrix(rows, rows);
        let b = generate_matrix(rows, rows);
        
        let a_view = MatrixView { data: &a.data, rows: a.rows, cols: a.cols, stride: a.cols };
        let b_view = MatrixView { data: &b.data, rows: b.rows, cols: b.cols, stride: b.cols };

        let start = std::time::Instant::now();
        let _c = divide_and_conquer(a_view, b_view);
        let elapsed = start.elapsed();

        total_us += elapsed.as_micros();
    }

    let avg_us = total_us / 20;

    println!(
        "DC avg: {}x{} -> {} µs",
        rows, rows, avg_us
    );

    if save {
        use std::io::Write;
        let mut file = std::fs::OpenOptions::new()
            .create(true)
            .append(true)
            .open(csv_path)
            .expect("Failed to open CSV");
    
        writeln!(file, "{},{}", rows, avg_us)
            .expect("Failed to write CSV");
    }
}
