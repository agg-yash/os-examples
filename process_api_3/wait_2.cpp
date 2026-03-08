#include <iostream>
#include <unistd.h>
#include <sys/wait.h>
#include <cstdlib>

using namespace std;

int main(int argc, char *argv[])
{
    cout << "hello world (pid:" << getpid() << ")" << endl;

    int rc = fork();

    if (rc < 0) { // fork failed
        cerr << "fork failed" << endl;
        exit(1);
    } 
    else if (rc == 0) { // child process
        cout << "hello, I am child (pid:" << getpid() << ")" << endl;
    } 
    else { // parent process
        int wc = wait(NULL);
        cout << "hello, I am parent of " << rc
             << " (wc:" << wc << ") (pid:" << getpid() << ")"
             << endl;
    }

    return 0;
}