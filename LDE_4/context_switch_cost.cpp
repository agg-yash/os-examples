#include <cerrno>
#include <cstdint>
#include <cstdlib>
#include <cstring>
#include <iostream>
#include <mach/mach.h>
#include <sys/time.h>
#include <sys/wait.h>
#include <unistd.h>

namespace {

constexpr long kDefaultRounds = 200000;
constexpr int kAffinityTag = 7;

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

void set_macos_affinity_tag(int tag) {
    thread_affinity_policy_data_t policy{};
    policy.affinity_tag = tag;
    const kern_return_t kr = thread_policy_set(
        mach_thread_self(),
        THREAD_AFFINITY_POLICY,
        reinterpret_cast<thread_policy_t>(&policy),
        THREAD_AFFINITY_POLICY_COUNT
    );
    if (kr != KERN_SUCCESS) {
        std::cerr << "warning: thread_policy_set failed (macOS affinity hint not applied)\n";
    }
}

} // namespace

int main(int argc, char* argv[]) {
    const long rounds = (argc >= 2) ? parse_long(argv[1], kDefaultRounds) : kDefaultRounds;

    int p2c[2];
    int c2p[2];
    if (pipe(p2c) != 0 || pipe(c2p) != 0) {
        std::cerr << "pipe() failed: " << std::strerror(errno) << "\n";
        return 1;
    }

    pid_t pid = fork();
    if (pid < 0) {
        std::cerr << "fork() failed: " << std::strerror(errno) << "\n";
        return 1;
    }

    char byte = 'x';

    if (pid == 0) {
        set_macos_affinity_tag(kAffinityTag);

        close(p2c[1]);
        close(c2p[0]);

        for (long i = 0; i < rounds; ++i) {
            if (read(p2c[0], &byte, 1) != 1) {
                _exit(2);
            }
            if (write(c2p[1], &byte, 1) != 1) {
                _exit(3);
            }
        }

        close(p2c[0]);
        close(c2p[1]);
        _exit(0);
    }

    set_macos_affinity_tag(kAffinityTag);

    close(p2c[0]);
    close(c2p[1]);

    const uint64_t start = now_us();
    for (long i = 0; i < rounds; ++i) {
        if (write(p2c[1], &byte, 1) != 1) {
            std::cerr << "parent write failed: " << std::strerror(errno) << "\n";
            return 1;
        }
        if (read(c2p[0], &byte, 1) != 1) {
            std::cerr << "parent read failed: " << std::strerror(errno) << "\n";
            return 1;
        }
    }
    const uint64_t end = now_us();

    close(p2c[1]);
    close(c2p[0]);

    int status = 0;
    waitpid(pid, &status, 0);
    if (!WIFEXITED(status) || WEXITSTATUS(status) != 0) {
        std::cerr << "child exited abnormally\n";
        return 1;
    }

    const double total_us = static_cast<double>(end - start);
    const double round_trip_us = total_us / static_cast<double>(rounds);
    // Each ping-pong round typically forces two context switches.
    const double context_switch_us = round_trip_us / 2.0;

    std::cout << "=== process ping-pong (pipes) ===\n";
    std::cout << "rounds: " << rounds << "\n";
    std::cout << "total time (us): " << total_us << "\n";
    std::cout << "avg round-trip (us): " << round_trip_us << "\n";
    std::cout << "estimated context switch cost (us/switch): " << context_switch_us << "\n";
    std::cout << "(estimate includes pipe read/write overhead)\n";

    return 0;
}

