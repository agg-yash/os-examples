#include <iostream>
#include <unistd.h>
#include <string.h>
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

        char *myargs[3];
        myargs[0] = strdup("wc");     // program to run
        myargs[1] = strdup("exec_3.cpp");   // file argument
        myargs[2] = NULL;             // argument list must end with NULL

        execvp(myargs[0], myargs);    // replace process with wc , also successfull call to exec never returns

        cout << "this shouldn't print out" << endl;
    } 
    else { // parent process
        int wc = wait(NULL); // wait for child process to finish , returns pid of child process
        cout << "hello, I am parent of " << rc
             << " (wc:" << wc << ") (pid:" << getpid() << ")"
             << endl;
    }

    return 0;
}