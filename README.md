# Comparison of Rust and Go Programming Languages Using Matrix Multiplication

This repository contains implementations and an experimental performance comparison of the **Rust** and **Go** programming languages using three matrix multiplication algorithms:

1. Standard iterative algorithm  
2. Divide-and-conquer algorithm  
3. Strassen’s algorithm  

The goal is to analyze execution time differences depending on the algorithm and the programming language.

---

## Running the Go implementation

Position inside /go directory and then:
```bash
go run .
```


---

## Running the Rust implementation

Position inside /rust/iterative directory and then:
```bash
cargo run --release [N]
```

N (optional) is the dimension of the square matrices.
