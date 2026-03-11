#include <iostream>
#include <unistd.h>
#include <string.h>
#include <fcntl.h>
#include <sys/wait.h>
#include <cstdlib>

using namespace std;

int main(int argc, char *argv[])
{
    int rc = fork();

    if (rc < 0) { // fork failed
        cerr << "fork failed" << endl;
        exit(1);
    }
    else if (rc == 0) { // child: redirect stdout to a file

        close(STDOUT_FILENO); // close stdout -> close(1)

        open("./redirect_4.output", O_CREAT | O_WRONLY | O_TRUNC, S_IRWXU);

        // now exec "wc"
        char *myargs[3];
        myargs[0] = strdup("wc");     // program
        myargs[1] = strdup("redirect_4.cpp");   // file argument
        myargs[2] = NULL;

        execvp(myargs[0], myargs);
    }
    else { // parent
        int wc = wait(NULL);
    }

    return 0;
}

// 0 , 1 --> redirect_4.output , 2
// 0 , 1 , 2