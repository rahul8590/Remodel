Remodel (CMPSCI 630)
=======

Remodel unlike make, it will use MD5 hashes to detect new content and provide, dependency build.

Remodel uses a different grammar than make. Dependencies can appear in any order. If you execute remodel with no arguments, it should start with the pseudo-target DEFAULT. Otherwise, the root is the argument to remodel, as in remodel foo.o. 

<code>
program ::= production*
production ::= target '<-' dependency (':' '"' command '"")
dependency ::= filename (',' filename)*
target ::= filename (',' filename)*
</code>

Here's an example that builds the program baz from two source files, foo.cpp and bar.cpp. 

<code>
DEFAULT <- baz
baz <- foo.o, bar.o: "g++ foo.o bar.o -o baz"
foo.o <- foo.cpp : "g++ -c foo.cpp -o foo.o"
bar.o <- bar.cpp: "g++ -c bar.cpp -o bar.o"
</code>

The dependencies are stored on disk in a special directory called .remodel/, so that remodel will not re-execute any commands unless a dependency has been violated.

