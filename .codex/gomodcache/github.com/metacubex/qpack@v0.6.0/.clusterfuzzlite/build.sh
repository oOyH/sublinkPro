#!/bin/bash -eu

export CXX="${CXX} -lresolv" # required by Go 1.20

compile_go_fuzzer github.com/metacubex/qpack/fuzzing Fuzz qpack_fuzzer
