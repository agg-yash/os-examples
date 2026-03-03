# intro_1

Small C++ programs meant to be run while you observe behavior of CPU and memory virtualization.

## Files

- **`cpu_vir.cpp`**: busy-waits for ~2 seconds, then prints the string you pass on the command line, forever. Useful to create a steady CPU load and watch it over time.
- **`mem_vir.cpp`**: prints PID and the addresses of a global, stack, and heap allocation, then loops forever. Useful for seeing a process memory layout at runtime.

## Build

From repo root:

```bash
g++ -O2 -std=c++17 -o intro_1/cpu_vir intro_1/cpu_vir.cpp
g++ -O2 -std=c++17 -o intro_1/mem_vir intro_1/mem_vir.cpp
```

## Run

### `cpu_vir`

```bash
./intro_1/cpu_vir "hello"
```

### `mem_vir`

```bash
./intro_1/mem_vir
```

