# process_2

`os_simulator.go` is a tiny scheduler/process-state simulator. It generates processes that either run CPU instructions or issue I/O, and prints a **time-based trace** showing how each process moves through:

- `READY`
- `RUNNING` (printed as `RUN:<opcode>` on the active PID)
- `BLOCKED` (waiting on I/O)
- `DONE`

The goal is to make it easy to see how scheduling policy choices change the timeline.

## Build / run

Run directly from repo root:

```bash
go run ./process_2/os_simulator.go -q 4:100,1:0 -x -m
```


## Flags (detailed)

### `-x`, `--execute`

If set, the simulator actually runs and prints the time trace table.

If **not** set, the program prints the generated instruction lists for each process (what *would* run) and prints the selected scheduling behaviors, then exits.

### `-m`, `--metrics`

Prints summary stats at the end (total time, CPU busy ticks, IO busy ticks). This is most useful with `-x/--execute`.

### `-q <list>`, `--queue <list>`

Defines processes to generate in the form:

`X1:Y1,X2:Y2,...`

Where:

- **X**: number of “instruction slots” to generate for that process
- **Y**: percent chance (0–100) that each slot is a CPU instruction

For each slot:

- with probability \(Y/100\): generate `cpu`
- otherwise: generate an `io` *and* an `io_done` (a 1-tick CPU instruction that “handles” I/O completion)

Examples:

- `-q 4:100` → one process: `cpu cpu cpu cpu`
- `-q 1:0` → one process: `io io_done`
- `-q 4:100,1:0` → two processes, CPU-bound + IO-bound

### `-g <programs>`, `--prog <programs>`

Defines explicit process programs (instead of random generation).

Format:

- Multiple processes: separate programs with `:`
- Within a program: comma-separated instructions

Instruction syntax:

- `cN` = N CPU ticks (e.g., `c3` → `cpu,cpu,cpu`)
- `i` = issue I/O (generates `io` then `io_done`)

Example (two processes):

```bash
go run ./process_2/os_simulator.go -g "c3,i,c1:i,c2" -x
```

Note: if `-g/--prog` is set, it takes precedence over `-q/--queue`.

### `-t <n>`, `--iotime <n>`

Sets how long each I/O takes (in simulator “time ticks”). Default is `5`.

When a process executes an `io` instruction, it becomes `BLOCKED` and an I/O completion event is scheduled `n` ticks later (plus the simulator’s internal offset).

### `-w <mode>`, `--switchwhen <mode>`

Controls when the CPU switches away from the current process:

- **`PREEMPT_ON_IO`** (default): switch when the running process **finishes** *or* **issues an I/O**
- **`SWITCH_ON_EXIT`**: switch only when the running process **finishes**

### `-e <mode>`, `--ioend <mode>`

Controls what happens when an I/O completes:

- **`WAKE_LATER`** (default): the process becomes `READY` and will run when it’s its turn
- **`WAKE_NOW`**: the process that completed I/O runs immediately (may preempt the current runner)

## Output columns

Each time tick prints:

- **Time**: current tick (a `*` means at least one I/O completed this tick)
- **PID columns**: state of each process (`READY`, `BLOCKED`, `DONE`, or `RUN:<opcode>`)
- **CPU**: `1` if an instruction ran this tick, blank otherwise
- **IOs**: number of I/Os in flight during this tick

## Example runs

### CPU-bound + IO-bound mix

```bash
go run ./process_2/os_simulator.go -q 4:100,1:0 -x -m
```

### Faster I/O

```bash
go run ./process_2/os_simulator.go -q 4:100,1:0 -t 2 -x -m
```

### Switch only when a process exits

```bash
go run ./process_2/os_simulator.go -q 6:50,6:50 -w SWITCH_ON_EXIT -x -m
```

### Run process immediately on I/O completion

```bash
go run ./process_2/os_simulator.go -q 8:40,8:40 -e WAKE_NOW -x -m
```

