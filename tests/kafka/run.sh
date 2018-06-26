#!/bin/bash

#
# @(#) Kafka Benchmark
#

export NDRX_BENCH_FILE=`pwd`/bench.txt
export NDRX_BENCH_CONFIGNAME="kafka_2.10-0.10.0.1, lin 4.10, 64bit, i5-4300U"
# clean up bech file..
> $NDRX_BENCH_FILE

pushd .

export PATH=$PATH:`pwd`/runtime/bin
# Debug settings:
export NDRX_CCONFIG=`pwd`/runtime/conf
export NDRX_APPHOME=`pwd`/runtime

rm runtime/log/* 2>/dev/null
cd runtime

CALLS=40000

#
# Generic exit function
#
function go_out {
    echo "Test exiting with: $1"
    xadmin killall rbtmqclt rbtmqsrv

    popd 2>/dev/null
    exit $1
}


echo "Starting server process..."
kfksrv &

SV_PID=$!

sleep 10

if ! kill -0 $SV_PID > /dev/null 2>&1; then
        echo "Kafka server not started! Is Kafka booted?" >&2
        go_out 1 
fi

echo "About to issue calls..."
kfkclt -num $CALLS
RET=$?

if [[ $RET != 0 ]]; then
	echo "kfkclt -num $CALLS failed"
	go_out 2
fi

go_out 0

