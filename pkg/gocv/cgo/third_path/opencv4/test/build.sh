#!/bin/bash

#complie
g++ test_imencode.cpp -std=c++11 -I../include -L../lib64 -lopencv_core -lopencv_highgui -lopencv_imgcodecs -o test_imencode

#run
#LD_LIBRARY_PATH=../lib64 ./test_imencode
