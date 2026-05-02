# LDE_4 - System Call and Context Switch Cost (macOS)

This exercise measures:

- the cost of a simple system call (`read(fd, nullptr, 0)`)
- the cost of a context switch using two processes ping-ponging through two pipes
- practical timer precision of `gettimeofday()`

## Files

- `syscall_cost.cpp`: probes `gettimeofday()` precision and then estimates null syscall cost.
- `context_switch_cost.cpp`: runs parent/child ping-pong over pipes and estimates context-switch cost.

## Build

From repo root:

```bash
clang++ -O2 -std=c++17 -Wall -Wextra -o LDE_4/syscall_cost LDE_4/syscall_cost.cpp
clang++ -O2 -std=c++17 -Wall -Wextra -o LDE_4/context_switch_cost LDE_4/context_switch_cost.cpp
```

(`g++` also works on most setups.)

## Run

### 1) Measure timer precision + null syscall cost

```bash
./LDE_4/syscall_cost
```

Optional args:

```bash
./LDE_4/syscall_cost <timer_samples> <syscall_iterations>
```

Example:

```bash
./LDE_4/syscall_cost 300000 8000000
```

### 2) Measure context-switch cost (pipe ping-pong)

```bash
./LDE_4/context_switch_cost
```

Optional arg:

```bash
./LDE_4/context_switch_cost <rounds>
```

Example:

```bash
./LDE_4/context_switch_cost 300000
```

## Notes for macOS

- Linux-style `sched_setaffinity()` is not available on macOS.
- The reported context-switch value is an estimate; it includes pipe communication overhead (read/write and synchronization), not just pure scheduler save/restore time.

