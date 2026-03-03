package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

const (
	SwitchOnIO  = "PREEMPT_ON_IO"
	SwitchOnEnd = "SWITCH_ON_EXIT"
)

const (
	ResumeLater     = "WAKE_LATER"
	ResumeImmediate = "WAKE_NOW"
)

const (
	StRunning = "RUNNING"
	StReady   = "READY"
	StDone    = "DONE"
	StBlocked = "BLOCKED"
)

const (
	kCode  = "code_"
	kPC    = "pc_"
	kPID   = "pid_"
	kState = "proc_state_"
)

const (
	OpCPU    = "cpu"
	OpIO     = "io"
	OpIODone = "io_done"
)

type scheduler struct {
	procInfo              map[int]map[string]any
	processSwitchBehavior string
	ioDoneBehavior        string
	ioLength              int
	currProc              int
	ioFinishTimes         map[int][]int
	rng                   *rand.Rand
}

func newScheduler(processSwitchBehavior, ioDoneBehavior string, ioLength int, rng *rand.Rand) *scheduler {
	return &scheduler{
		procInfo:              map[int]map[string]any{},
		processSwitchBehavior: processSwitchBehavior,
		ioDoneBehavior:        ioDoneBehavior,
		ioLength:              ioLength,
		rng:                   rng,
	}
}

func (s *scheduler) newProcess() int {
	procID := len(s.procInfo)
	s.procInfo[procID] = map[string]any{}
	s.procInfo[procID][kPC] = 0
	s.procInfo[procID][kPID] = procID
	s.procInfo[procID][kCode] = []string{}
	s.procInfo[procID][kState] = StReady
	return procID
}

func (s *scheduler) loadProgram(program string) {
	procID := s.newProcess()
	for _, part := range strings.Split(program, ",") {
		if part == "" {
			fatalf("bad opcode %q (should be c or i)", part)
		}
		opcode := part[0]
		switch opcode {
		case 'c':
			num, err := strconv.Atoi(part[1:])
			if err != nil {
				fatalf("bad compute count in %q", part)
			}
			code := s.procInfo[procID][kCode].([]string)
			for i := 0; i < num; i++ {
				code = append(code, OpCPU)
			}
			s.procInfo[procID][kCode] = code
		case 'i':
			code := s.procInfo[procID][kCode].([]string)
			code = append(code, OpIO)
			code = append(code, OpIODone)
			s.procInfo[procID][kCode] = code
		default:
			fatalf("bad opcode %c (should be c or i)", opcode)
		}
	}
}

func (s *scheduler) load(programDescription string) {
	procID := s.newProcess()
	tmp := strings.Split(programDescription, ":")
	if len(tmp) != 2 {
		fmt.Printf("Bad description (%s): Must be number <x:y>\n", programDescription)
		fmt.Printf("  where X is the number of instructions\n")
		fmt.Printf("  and Y is the percent change that an instruction is CPU not IO\n")
		os.Exit(1)
	}

	numInstructions, err := strconv.Atoi(tmp[0])
	if err != nil {
		fatalf("Bad X in %q", programDescription)
	}
	chanceCPUPercent, err := strconv.ParseFloat(tmp[1], 64)
	if err != nil {
		fatalf("Bad Y in %q", programDescription)
	}
	chanceCPU := chanceCPUPercent / 100.0

	code := s.procInfo[procID][kCode].([]string)
	for i := 0; i < numInstructions; i++ {
		if s.rng.Float64() < chanceCPU {
			code = append(code, OpCPU)
		} else {
			code = append(code, OpIO)
			code = append(code, OpIODone)
		}
	}
	s.procInfo[procID][kCode] = code
}

func (s *scheduler) moveToReady(expected string, pid int) {
	if pid == -1 {
		pid = s.currProc
	}
	if s.procInfo[pid][kState].(string) != expected {
		panic(fmt.Sprintf("assert failed: pid %d expected %s, got %s", pid, expected, s.procInfo[pid][kState].(string)))
	}
	s.procInfo[pid][kState] = StReady
}

func (s *scheduler) moveToWait(expected string) {
	if s.procInfo[s.currProc][kState].(string) != expected {
		panic(fmt.Sprintf("assert failed: pid %d expected %s, got %s", s.currProc, expected, s.procInfo[s.currProc][kState].(string)))
	}
	s.procInfo[s.currProc][kState] = StBlocked
}

func (s *scheduler) moveToRunning(expected string) {
	if s.procInfo[s.currProc][kState].(string) != expected {
		panic(fmt.Sprintf("assert failed: pid %d expected %s, got %s", s.currProc, expected, s.procInfo[s.currProc][kState].(string)))
	}
	s.procInfo[s.currProc][kState] = StRunning
}

func (s *scheduler) moveToDone(expected string) {
	if s.procInfo[s.currProc][kState].(string) != expected {
		panic(fmt.Sprintf("assert failed: pid %d expected %s, got %s", s.currProc, expected, s.procInfo[s.currProc][kState].(string)))
	}
	s.procInfo[s.currProc][kState] = StDone
}

func (s *scheduler) nextProc(pid int) {
	if pid != -1 {
		s.currProc = pid
		s.moveToRunning(StReady)
		return
	}
	for p := s.currProc + 1; p < len(s.procInfo); p++ {
		if s.procInfo[p][kState].(string) == StReady {
			s.currProc = p
			s.moveToRunning(StReady)
			return
		}
	}
	for p := 0; p <= s.currProc; p++ {
		if s.procInfo[p][kState].(string) == StReady {
			s.currProc = p
			s.moveToRunning(StReady)
			return
		}
	}
}

func (s *scheduler) getNumProcesses() int {
	return len(s.procInfo)
}

func (s *scheduler) getNumInstructions(pid int) int {
	return len(s.procInfo[pid][kCode].([]string))
}

func (s *scheduler) getInstruction(pid, index int) string {
	return s.procInfo[pid][kCode].([]string)[index]
}

func (s *scheduler) getNumActive() int {
	numActive := 0
	for pid := 0; pid < len(s.procInfo); pid++ {
		if s.procInfo[pid][kState].(string) != StDone {
			numActive++
		}
	}
	return numActive
}

func (s *scheduler) getNumRunnable() int {
	numActive := 0
	for pid := 0; pid < len(s.procInfo); pid++ {
		st := s.procInfo[pid][kState].(string)
		if st == StReady || st == StRunning {
			numActive++
		}
	}
	return numActive
}

func (s *scheduler) getIOsInFlight(currentTime int) int {
	numInFlight := 0
	for pid := 0; pid < len(s.procInfo); pid++ {
		for _, t := range s.ioFinishTimes[pid] {
			if t > currentTime {
				numInFlight++
			}
		}
	}
	return numInFlight
}

func (s *scheduler) checkIfDone() {
	if len(s.procInfo[s.currProc][kCode].([]string)) == 0 {
		if s.procInfo[s.currProc][kState].(string) == StRunning {
			s.moveToDone(StRunning)
			s.nextProc(-1)
		}
	}
}

func (s *scheduler) run() (cpuBusy, ioBusy, clockTick int) {
	clockTick = 0
	if len(s.procInfo) == 0 {
		return 0, 0, 0
	}

	s.ioFinishTimes = map[int][]int{}
	for pid := 0; pid < len(s.procInfo); pid++ {
		s.ioFinishTimes[pid] = []int{}
	}

	s.currProc = 0
	s.moveToRunning(StReady)

	fmt.Printf("%s", "Time")
	for pid := 0; pid < len(s.procInfo); pid++ {
		fmt.Printf("%14s", fmt.Sprintf("PID:%2d", pid))
	}
	fmt.Printf("%14s", "CPU")
	fmt.Printf("%14s", "IOs")
	fmt.Printf("\n")

	ioBusy = 0
	cpuBusy = 0

	for s.getNumActive() > 0 {
		clockTick++

		ioDone := false
		for pid := 0; pid < len(s.procInfo); pid++ {
			if containsInt(s.ioFinishTimes[pid], clockTick) {
				ioDone = true
				s.moveToReady(StBlocked, pid)
				if s.ioDoneBehavior == ResumeImmediate {
					if s.currProc != pid {
						if s.procInfo[s.currProc][kState].(string) == StRunning {
							s.moveToReady(StRunning, -1)
						}
					}
					s.nextProc(pid)
				} else {
					if s.processSwitchBehavior == SwitchOnEnd && s.getNumRunnable() > 1 {
						s.nextProc(pid)
					}
					if s.getNumRunnable() == 1 {
						s.nextProc(pid)
					}
				}
				s.checkIfDone()
			}
		}

		instructionToExecute := ""
		if s.procInfo[s.currProc][kState].(string) == StRunning &&
			len(s.procInfo[s.currProc][kCode].([]string)) > 0 {
			code := s.procInfo[s.currProc][kCode].([]string)
			instructionToExecute = code[0]
			s.procInfo[s.currProc][kCode] = code[1:]
			cpuBusy++
		}

		if ioDone {
			fmt.Printf("%3d*", clockTick)
		} else {
			fmt.Printf("%3d ", clockTick)
		}
		for pid := 0; pid < len(s.procInfo); pid++ {
			if pid == s.currProc && instructionToExecute != "" {
				fmt.Printf("%14s", "RUN:"+instructionToExecute)
			} else {
				fmt.Printf("%14s", s.procInfo[pid][kState].(string))
			}
		}

		if instructionToExecute == "" {
			fmt.Printf("%14s", " ")
		} else {
			fmt.Printf("%14s", "1")
		}

		numOutstanding := s.getIOsInFlight(clockTick)
		if numOutstanding > 0 {
			fmt.Printf("%14s", strconv.Itoa(numOutstanding))
			ioBusy++
		} else {
			fmt.Printf("%10s", " ")
		}
		fmt.Printf("\n")

		if instructionToExecute == OpIO {
			s.moveToWait(StRunning)
			s.ioFinishTimes[s.currProc] = append(s.ioFinishTimes[s.currProc], clockTick+s.ioLength+1)
			if s.processSwitchBehavior == SwitchOnIO {
				s.nextProc(-1)
			}
		}

		s.checkIfDone()
	}

	return cpuBusy, ioBusy, clockTick
}

func containsInt(xs []int, v int) bool {
	for _, x := range xs {
		if x == v {
			return true
		}
	}
	return false
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func main() {
	prog := flag.String("g", "", "explicit programs, colon-separated (e.g., c7,i:c3,i)")
	progLong := flag.String("prog", "", "explicit programs (same as -g)")

	procList := flag.String("q", "", "process list X1:Y1,X2:Y2,... (X instructions, Y% CPU vs IO)")
	procListLong := flag.String("queue", "", "process list (same as -q)")

	ioLen := flag.Int("t", 5, "how long an IO takes")
	ioLenLong := flag.Int("iotime", 5, "how long an IO takes (same as -t)")

	sw := flag.String("w", SwitchOnIO, "when to switch: PREEMPT_ON_IO, SWITCH_ON_EXIT")
	swLong := flag.String("switchwhen", SwitchOnIO, "when to switch (same as -w)")

	ioDone := flag.String("e", ResumeLater, "when IO ends: WAKE_LATER, WAKE_NOW")
	ioDoneLong := flag.String("ioend", ResumeLater, "when IO ends (same as -e)")

	doRun := flag.Bool("x", false, "execute and print the trace")
	doRunLong := flag.Bool("execute", false, "execute (same as -x)")

	printStats := flag.Bool("m", false, "print statistics at end (useful with -x)")
	printStatsLong := flag.Bool("metrics", false, "print statistics at end (same as -m)")

	flag.Parse()

	seed := 0

	program := *prog
	if *progLong != "" {
		program = *progLong
	}

	processList := *procList
	if *procListLong != "" {
		processList = *procListLong
	}

	iolength := *ioLen
	if *ioLenLong != 5 {
		iolength = *ioLenLong
	}

	switchBehavior := *sw
	if *swLong != SwitchOnIO {
		switchBehavior = *swLong
	}

	ioDoneBehavior := *ioDone
	if *ioDoneLong != ResumeLater {
		ioDoneBehavior = *ioDoneLong
	}

	solve := *doRun || *doRunLong
	stats := *printStats || *printStatsLong

	if switchBehavior != SwitchOnIO && switchBehavior != SwitchOnEnd {
		fatalf("bad switch behavior %q (must be %s or %s)", switchBehavior, SwitchOnIO, SwitchOnEnd)
	}
	if ioDoneBehavior != ResumeImmediate && ioDoneBehavior != ResumeLater {
		fatalf("bad io-done behavior %q (must be %s or %s)", ioDoneBehavior, ResumeLater, ResumeImmediate)
	}
	if iolength < 0 {
		fatalf("iolength must be >= 0")
	}

	rng := rand.New(rand.NewSource(int64(seed)))
	s := newScheduler(switchBehavior, ioDoneBehavior, iolength, rng)

	if program != "" {
		for _, p := range strings.Split(program, ":") {
			s.loadProgram(p)
		}
	} else {
		for _, p := range strings.Split(processList, ",") {
			s.load(p)
		}
	}

	if !solve {
		fmt.Printf("Produce a trace of what would happen when you run these processes:\n")
		for pid := 0; pid < s.getNumProcesses(); pid++ {
			fmt.Printf("Process %d\n", pid)
			for inst := 0; inst < s.getNumInstructions(pid); inst++ {
				fmt.Printf("  %s\n", s.getInstruction(pid, inst))
			}
			fmt.Printf("\n")
		}
		fmt.Printf("Important behaviors:\n")
		fmt.Printf("  System will switch when ")
		if switchBehavior == SwitchOnIO {
			fmt.Printf("the current process is FINISHED or ISSUES AN IO\n")
		} else {
			fmt.Printf("the current process is FINISHED\n")
		}
		fmt.Printf("  After IOs, the process issuing the IO will ")
		if ioDoneBehavior == ResumeImmediate {
			fmt.Printf("run IMMEDIATELY\n")
		} else {
			fmt.Printf("run LATER (when it is its turn)\n")
		}
		fmt.Printf("\n")
		os.Exit(0)
	}

	cpuBusy, ioBusy, clockTick := s.run()
	if stats {
		fmt.Printf("\n")
		fmt.Printf("Stats: Total Time %d\n", clockTick)
		fmt.Printf("Stats: CPU Busy %d (%.2f%%)\n", cpuBusy, 100.0*float64(cpuBusy)/float64(clockTick))
		fmt.Printf("Stats: IO Busy  %d (%.2f%%)\n", ioBusy, 100.0*float64(ioBusy)/float64(clockTick))
		fmt.Printf("\n")
	}
}
