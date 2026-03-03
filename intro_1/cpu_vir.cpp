#include <iostream>
#include <cstdlib>
#include <string>
#include <chrono>

using namespace std;

void Wait(int seconds) {
    auto start = chrono::steady_clock::now();
    while (true) {
        auto now = chrono::steady_clock::now();
        auto elapsed =
            chrono::duration_cast<chrono::seconds>(now - start).count();
        if (elapsed >= seconds) {
            break;
        }
    }
}

int main(int argc, char* argv[]) {

    string str = argv[1];

    while (true) {
        Wait(2);
        cout << str << endl;
    }

    return 0;
}
