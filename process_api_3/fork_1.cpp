#include <iostream>
#include <unistd.h>
#include <sys/types.h>
#include <cstdlib>

using namespace std;


int main(int argc, char *argv[])
{
    cout << "hello world (pid:" << getpid() << ")" << endl;

    int rc = fork(); // fork system call

    if (rc < 0) { // fork failed
        cerr << "fork failed" << endl;
        exit(1);
    } 
    else if (rc == 0) { // child process : child process will always have rc 0
        cout << "hello, I am child (pid:" << getpid() << ")" << std::endl;
    } 
    else { // parent process : parent process will always have rc as pid of child process
        cout << "hello, I am parent of " << rc << " (pid:" << getpid() << ")" << endl;
    }

    return 0;
}
