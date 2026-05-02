#include <cerrno>
#include <cstdint>
#include <cstdlib>
#include <cstring>
#include <iostream>
#include <sys/time.h>
#include <unistd.h>

namespace {

constexpr long kDefaultTimerSamples = 200000;
constexpr long kDefaultSyscallIters = 5000000;

inline uint64_t now_us() {
    timeval tv{};
    gettimeofday(&tv, nullptr);
    return static_cast<uint64_t>(tv.tv_sec) * 1000000ULL + static_cast<uint64_t>(tv.tv_usec);
}

long parse_long(const char* s, long fallback) {
    if (s == nullptr) {
        return fallback;
    }
    char* end = nullptr;
    errno = 0;
    const long v = std::strtol(s, &end, 10);
    if (errno != 0 || end == s || *end != '\0' || v <= 0) {
        return fallback;
    }
    return v;
}

} // namespace

int main(int argc, char* argv[]) {
    const long timer_samples = (argc >= 2) ? parse_long(argv[1], kDefaultTimerSamples) : kDefaultTimerSamples;
    const long syscall_iters = (argc >= 3) ? parse_long(argv[2], kDefaultSyscallIters) : kDefaultSyscallIters;

    uint64_t nonzero_delta_sum = 0;
    uint64_t min_nonzero_delta = UINT64_MAX;
    long nonzero_count = 0;

    for (long i = 0; i < timer_samples; ++i) {
        const uint64_t t1 = now_us();
        const uint64_t t2 = now_us();
        const uint64_t d = t2 - t1;
        if (d > 0) {
            nonzero_delta_sum += d;
            min_nonzero_delta = (d < min_nonzero_delta) ? d : min_nonzero_delta;
            ++nonzero_count;
        }
    }

    const uint64_t start = now_us();
    for (long i = 0; i < syscall_iters; ++i) {
        // A 0-byte read still enters the kernel and returns immediately.
        (void) read(STDIN_FILENO, nullptr, 0);
    }
    const uint64_t end = now_us();

    const double total_us = static_cast<double>(end - start);
    const double syscall_us = total_us / static_cast<double>(syscall_iters);

    std::cout << "=== gettimeofday() precision probe ===\n";
    std::cout << "samples: " << timer_samples << "\n";
    if (nonzero_count == 0) {
        std::cout << "all back-to-back deltas were 0 us (timer granularity is coarse at this scale)\n";
    } else {
        const double avg_nonzero = static_cast<double>(nonzero_delta_sum) / static_cast<double>(nonzero_count);
        std::cout << "non-zero deltas observed: " << nonzero_count << "\n";
        std::cout << "min non-zero delta (us): " << min_nonzero_delta << "\n";
        std::cout << "avg non-zero delta (us): " << avg_nonzero << "\n";
    }

    std::cout << "\n=== null syscall cost (0-byte read) ===\n";
    std::cout << "iterations: " << syscall_iters << "\n";
    std::cout << "total time (us): " << total_us << "\n";
    std::cout << "estimated syscall cost (us/call): " << syscall_us << "\n";

    return 0;
}

