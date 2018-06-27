#!/bin/bash

#
# @(#) Kafka Benchmark
#

export NDRX_BENCH_FILE=`pwd`/bench1w.txt
export NDRX_BENCH_CONFIGNAME="apache-activemq-5.15.4, lin 4.10, 64bit, i5-4300U"

# clean up bech file..
> $NDRX_BENCH_FILE

pushd .

export PATH=$PATH:`pwd`/runtime/bin
# Debug settings:
export NDRX_CCONFIG=`pwd`/runtime/conf
export NDRX_APPHOME=`pwd`/runtime

rm runtime/log/* 2>/dev/null
cd runtime

CALLS=400000

#
# Generic exit function
#
function go_out {
    echo "Test exiting with: $1"
    xadmin killall amq

    popd 2>/dev/null
    exit $1
}


echo "Starting server process..."
amqsrvoneway &

SV_PID=$!

sleep 10

if ! kill -0 $SV_PID > /dev/null 2>&1; then
        echo "ActiveMQ server not started! Is ActiveMQ booted?" >&2
        go_out 1 
fi

amqclt -num $CALLS -oneway
RET=$?

if [[ $RET != 0 ]]; then
	echo "amqclt -num $CALLS -oneway failed"
	go_out 2
fi


RESULT=""

while [ "X$RESULT" == "X" ]; do
        sleep 10
        echo "Checking for completion..."
        RESULT=`grep ' 4032 ' $NDRX_BENCH_FILE`
done

go_out 0

