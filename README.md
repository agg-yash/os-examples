# OS (learning / experiments)

This repo is a collection of OS-related experiments and exampples (CPU behavior, memory layout, and basic process scheduling/state tracing).

 ## Learning videos

I explain these experiments on my YouTube channel: [Bare Metal](https://www.youtube.com/@BareMetal-1)

## Contents

- **`intro_1/`**: tiny C++ programs to observe CPU + memory virtualization
- **`process_2/`**: a simple process/IO scheduling simulator (prints a time-based trace)

## Quick start

### `intro_1/` (C++)

Build (from repo root):

```bash
g++ -O2 -std=c++17 -o intro_1/cpu_vir intro_1/cpu_vir.cpp
g++ -O2 -std=c++17 -o intro_1/mem_vir intro_1/mem_vir.cpp
```

Run:

```bash
./intro_1/cpu_vir "hello"
./intro_1/mem_vir
```

### `process_2/` (Go)

Run directly:

```bash
go run ./process_2/os_simulator.go -q 4:100,1:0 -x -m
```

See `process_2/README.md` for the full flag reference and examples.

