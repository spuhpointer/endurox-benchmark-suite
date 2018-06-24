#!/bin/bash

#
# @(#) Generate the chart...
#

export NDRX_BENCH_FILE=`pwd`/vs1w.txt
echo '"Configuration" "MsgSize" "CallsPerSec"' > $NDRX_BENCH_FILE 

find . -name bench1w.txt | xargs -i cat {} | grep -v CallsPerSec >> $NDRX_BENCH_FILE

#
# Generate the chart
#
export NDRX_BENCH_TITLE="Middleware benchmarks - Send only"
export NDRX_BENCH_X_LABEL="Msg Size (bytes)"
export NDRX_BENCH_Y_LABEL="Msgs Per Second"
export NDRX_BENCH_OUTFILE="vs1w.png"
R -f genchart.r


