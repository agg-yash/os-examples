## process_api_3

Small C++ programs to explore basic UNIX process APIs:

- `fork(2)`
- `wait(2)` / `waitpid(2)`
- `execve(2)` (via the `exec*` family)
- simple output redirection with `dup2(2)`

Each file is a tiny, focused example you can compile and run while watching process trees, exit statuses, and file descriptors.

### Files

- `fork_1.cpp`: create a child process and print PIDs from both parent and child to see how `fork()` duplicates execution.
- `wait_2.cpp`: demonstrate how a parent can wait for a child to finish and inspect its exit status.
- `exec_3.cpp`: show how `exec()` replaces the current process image with a new program.
- `redirect_4.cpp`: redirect standard output to a file using `dup2()` and then run code that writes to stdout.

### Build

From repo root:

```bash
g++ -O2 -std=c++17 -o process_api_3/fork_1 process_api_3/fork_1.cpp
g++ -O2 -std=c++17 -o process_api_3/wait_2 process_api_3/wait_2.cpp
g++ -O2 -std=c++17 -o process_api_3/exec_3 process_api_3/exec_3.cpp
g++ -O2 -std=c++17 -o process_api_3/redirect_4 process_api_3/redirect_4.cpp
```

### Run

```bash
./process_api_3/fork_1
./process_api_3/wait_2
./process_api_3/exec_3
./process_api_3/redirect_4
```

