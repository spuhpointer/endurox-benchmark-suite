#!/bin/bash

#
# @(#) Enduro/X Benchmark
#

export NDRX_BENCH_FILE=`pwd`/bench.txt
export NDRX_BENCH_CONFIGNAME="Enduro/X 5.4 beta, on Linux 4.10, i5-4300U, Golang"

# clean up bech file..
> $NDRX_BENCH_FILE


pushd .

rm runtime/log/* 2>/dev/null
cd runtime

CALLS=400000

#
# Generic exit function
#
function go_out {
    echo "Test exiting with: $1"
    xadmin stop -y
    xadmin down -y

    popd 2>/dev/null
    exit $1
}

#
# Generate the runtime
#
xadmin provision -d

cd conf
. settest1

cd ..

# Start the system
xadmin start -y

exclt -num $CALLS -retry 2
RET=$?

if [[ $RET != 0 ]]; then
	echo "exclt -num $CALLS failed"
	go_out 1
fi

go_out 0

