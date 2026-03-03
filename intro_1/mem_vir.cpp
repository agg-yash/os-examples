#include <unistd.h>
#include <iostream>
#include <chrono>

using namespace std;

void Wait(int seconds) {
    auto start = chrono::steady_clock::now();
    while (chrono::duration_cast<chrono::seconds>(
           chrono::steady_clock::now() - start).count() < seconds) {}
}

int global_var = 42;

int main() {
    int stack_var = 10;
    int* heap_var = new int(20);

    cout << "PID: " << getpid() << endl;
    cout << "Address of global: " << &global_var << endl;
    cout << "Address of stack:  " << &stack_var << endl;
    cout << "Address of heap:   " << heap_var << endl;

    while (true) {
       Wait(5); 
    }
}
