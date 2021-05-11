Installation Procedure for Downloading Database
How to Install Golang (If Not Downloaded) to Run Source Code
Mac:
Change to the directory where you intend to install Go. Then clone the repository and check out the latest release tag 
Follow steps in this link if git not installed: https://git-scm.com/book/en/v2/Getting-Started-Installing-Git 

$ git clone https://go.googlesource.com/go goroot
$ cd goroot
$ git checkout go1.16.4     <-look up latest version of Go 
Build the Go distribution
$ cd src
$ ./all.bash

Open a new shell session. Now prepend the bin directory from the output above to your PATH so that you are using your custom go binary by default.
Mac:
$ export PATH=<Your path to go src>/src/go/bin:$PATH
Example: $ export PATH=/Users/carolynvs/src/go/bin:$PATH
Windows:
Download Go for Windows
https://golang.org/dl/go1.16.4.darwin-amd64.pkg <- click this link to automatically download Go onto windows
Verify that it’s been installed”
$ go version

How to Compile and Run Source Code
Download attached Database_ProjectCPSC332 folder 
or….. GitHub Optional Process: Clone repo code from Github using this link https://github.com/clairerhoda/Database_CPSC332.git
In a terminal window, go into the folder where Database_ProjectCPSC332 was downloaded
Use command go build to compile the code
Use command go run main.go to run code on localhost:8080 server
Use Postman to make calls to the Database (path names can be found within Go files) Example: http://localhost:8080/getDepartments
