#!/bin/bash

#
# @(#) Generate the chart...
#

export NDRX_BENCH_FILE=`pwd`/vs.txt
echo '"Configuration" "MsgSize" "CallsPerSec"' > $NDRX_BENCH_FILE 

find . -name bench.txt | xargs -i cat {} | grep -v CallsPerSec >> $NDRX_BENCH_FILE

#
# Generate the chart
#
export NDRX_BENCH_TITLE="Middleware benchmarks"
export NDRX_BENCH_X_LABEL="Msg Size (bytes)"
export NDRX_BENCH_Y_LABEL="Calls Per Second"
export NDRX_BENCH_OUTFILE="vs.png"
R -f genchart.r


