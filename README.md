Remodel (CMPSCI 630)
=======

Remodel unlike make, it will use MD5 hashes to detect new content and provide, dependency build.

Remodel uses a different grammar than make. Dependencies can appear in any order. If you execute remodel with no arguments, it should start with the pseudo-target DEFAULT. Otherwise, the root is the argument to remodel, as in remodel foo.o. 

    program ::= production*
    production ::= target '<-' dependency (':' '"' command '"")
    dependency ::= filename (',' filename)*
    target ::= filename (',' filename)*
	

Here's an example that builds the program baz from two source files, foo.cpp and bar.cpp. 

    DEFAULT <- baz
    baz <- foo.o, bar.o: "g++ foo.o bar.o -o baz"
    foo.o <- foo.cpp : "g++ -c foo.cpp -o foo.o"
    bar.o <- bar.cpp: "g++ -c bar.cpp -o bar.o"
	

The dependencies are stored on disk in a special directory called .remodel/, so that remodel will not re-execute any commands unless a dependency has been violated.


Usage
------

Remodel is pretty simple to use. Have all the build dependencies writted to config file
named "config" (without the quotes ofcourse :) . If no parameters are passed to it, default will be
choosen as root.

```bash
 $ go build main.go
 $ ./main <optional_root_name>
```

The following the go version

```bash
  $ go version
  go version go1
```


> Note: The program assumes all the files are present in the same directory as main.go
> Currently, it does not support files present in sub-directories :(


Test Cases
----------

1. Check for 1st time Scan 
2. Check Intermediate Builds 
3. Builds based on Checkpoints 
4. Check For Cyclic Dependencies
5. Handling of Missing source files / Intermediate object files


More Page Links
---------------
[Code Walkthrough](Remodel/wiki/Code-Walkthrough)