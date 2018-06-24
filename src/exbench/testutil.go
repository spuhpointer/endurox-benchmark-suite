package exbench

/*
 * @brief Test util package
 *
 * @file testutil.go
 *
 */

/*
#cgo pkg-config: atmisrvinteg

#include <ndebug.h>
#include <string.h>
#include <stdlib.h>
#include <nstdutil.h>


*/
import "C"

import (
	"flag"

	atmi "github.com/endurox-dev/endurox-go"
)

/**
 * SUCCEED/FAIL flags
 */
const (
	/** succeed status */
	SUCCEED = 0
	/** fail status */
	FAIL = -1

	/** number of messages to war up the system */
	WARMUP = 2
)

/**
 * Write benchmark results
 * @param[in] message size at benchmark (bytes)
 * @param[in] callspersec number of calls per second
 * @return 0 - succeed; -1 - fail
 */
func Ndrx_bench_write_stats(msgsize float64, callspersec float64) int {
	return int(C.ndrx_bench_write_stats(C.double(msgsize), C.double(callspersec)))
}

/**
 * First buffer prepare for given size
 * @param size[in] array size for testing
 * @param mod[in] modulus for the counter, for different tests
 * @return allocated test buffer
 */
func Ndrx_bench_get_buffer(size int, mod int) []byte {

	ret := make([]byte, size)

	for i := 0; i < size; i++ {
		ret[i] = byte(255 - i%mod)
	}

	return ret
}

/**
 * Verify the first buffer for the given size
 * @param buf[in] buffer to test
 * @param size[in] buffer size (from different source to check the actual len)
 * @param mod[in] value modulus (for different tests)
 * @return -1 fail, 0 succeed
 */
func Ndrx_bench_verify_buffer(ctx *atmi.ATMICtx, buf []byte, size int, mod int, msg string) int {

	if size <= 0 {
		ctx.TpLogError("TESTERROR! Invalid buffer size received: %d", size)
		return FAIL
	}

	for i := 0; i < size; i++ {
		expected := byte(255 - i%mod)

		if expected != buf[i] {
			ctx.TpLogError("%s at position %d modulus %d expected: %d got %d",
				msg, i, mod, expected, buf[i])

			ctx.TpLogDump(atmi.LOG_ERROR, "Invalid buffer", buf, len(buf))
			return FAIL
		}
	}

	return SUCCEED
}

/**
 * We have received message from client (at server side)
 * @param ctx ATMI Context (for logging)
 * @param correl Call Correlator
 * @param buf buffer received
 * @return status code 0 ok, -1 fail
 * @return byte array to send away if OK
 */
func Ndrx_bench_svmain(ctx *atmi.ATMICtx, correl int64, buf []byte) (int, []byte) {

	/* we got buffer, lets very it... */

	if ret := Ndrx_bench_verify_buffer(ctx, buf, len(buf), 255,
		"TESTERROR! Invalid data from client!"); ret != SUCCEED {
		ctx.TpLogError("TESTERROR! Invalid data from client, received by server")
		return FAIL, nil
	}

	/* generate reply buffer */
	retbuf := Ndrx_bench_get_buffer(len(buf), 254)

	return SUCCEED, retbuf
}

/** Buffer len (last tested) */
var M_bufLen int = 1

/** Number of calls received */
var M_nrcalls int = 0

/** Stopwatch for server tests */
var M_w StopWatch

/**
 * One way server function (measure that calls).
 * Note that measurements will include the the time for receival of the next buffer
 * size. Due to larg quantity of the messages, it is expected that it will not cause
 * significant differentces in  the results.
 * @param ctx ATMI Context for logging
 * @param correl call correlator (optional)
 * @param buf Buffer received to verify
 * @return 0 succeed, -1 fail
 */
func Ndrx_bench_svmain_oneway(ctx *atmi.ATMICtx, correl int64, buf []byte) int {

	/* verify incoming buffer */

	rcv_len := len(buf)

	if ret := Ndrx_bench_verify_buffer(ctx, buf, rcv_len, 255,
		"TESTERROR! Invalid data from client!"); ret != SUCCEED {
		ctx.TpLogError("TESTERROR! Invalid data from client, received by server")
		return FAIL
	}

	M_nrcalls++

	/* start the measurements if buflen is other than 1 */
	if M_bufLen != rcv_len {

		if M_bufLen != 1 {
			/* we got next call, lets plot results */
			if ret := Ndrx_bench_write_stats(float64(M_bufLen),
				float64(M_nrcalls)/float64(M_w.GetDetlaSec())); SUCCEED != ret {
				ctx.TpLogError("Failed to write benchmark stats!")
				return FAIL
			}
		}
		M_w.Reset()
		M_bufLen = rcv_len
		M_nrcalls = 1
	}

	return SUCCEED
}

/**
 * request callback function
 * @param ctx[in] ATMI context (for logging, etc...)
 * @param buf[in] Byte buffer to send
 * @param correl[in] Call correlator (used by systems where needed to match req with rsp)
 * @return status code 0 = succeed, -1 = FAIL, return buffer.
 */
type Ndrx_Bench_requestCB func(ctx *atmi.ATMICtx, correl int64, buf []byte, oneway bool) (int, []byte)

/**
 * Benchmark main
 * @param threadid[in] logical thread id used for correlators...
 * @param request[in] callback function for sending the data to server and
 *   receiving response back
 * @return 0 = succeed, -1 fail
 */
func Ndrx_bench_clmain(ctx *atmi.ATMICtx, threadid int, request Ndrx_Bench_requestCB) int {

	nrrequestsPtr := flag.Int("num", 500000, "Number of requests")
	retryPtr := flag.Int("retry", 1, "Number of retries")
	onewayPtr := flag.Bool("oneway", false, "a bool")

	flag.Parse()

	size := 0
	retriesDone := 1

	nrrequests := *nrrequestsPtr
	retries := *retryPtr
	oneway := *onewayPtr

	ctx.TpLogInfo("Number of request per message size: %d, retries: %d, oneway: %t",
		nrrequests, retries, oneway)

	/* Correlator shall be built as thread id + call number (we need some offset
	 * for negative steps)

	 */
	for i := -1 * WARMUP; i < 128; i++ {

		/* have some time for warmup */
		if i <= 0 {
			size = 1
		} else {
			size = i * 32
		}

		ctx.TpLogInfo("Benchmarking step: %d with size of %d bytes (nr req: %d)",
			i, size, nrrequests)

		/* we shall loop over the given count, let say nrrequests / 100K requests
		 * and shall start the stopwatch
		 */
		var w StopWatch

		w.Reset()
		requests_to_do := nrrequests / (WARMUP + i + 1)

		for req := 0; req < requests_to_do; req++ {

			buf := Ndrx_bench_get_buffer(size, 255)

			/* build up correlator... */
			var correl int64

			correl = int64(threadid)

			correl <<= 8

			correl |= WARMUP + int64(i)

			correl <<= 40

			correl |= int64(req)
		restart:
			res, retbuf := request(ctx, correl, buf, oneway)

			if res != SUCCEED {
				ctx.TpLogError("Server failure on call=%d size=%d iter=%d correl=%d",
					req, size, i, correl)

				if retriesDone < retries {
					retriesDone++
					ctx.TpLogWarn("Got reply error ... restarting %d", retriesDone)
					goto restart
				}

				return FAIL

			}

			/* verify only if it is RPC.. */
			if !oneway {
				if ret := Ndrx_bench_verify_buffer(ctx, retbuf, size, 254, "Error in reply"); ret != SUCCEED {

					ctx.TpLogError("Failed to verify reply buffer on call=%d size=%d iter=%d correl=%d",
						req, size, i, correl)
					return FAIL
				}
			}
		}

		/* ready to plot the results... */

		if i >= 0 && !oneway {
			if ret := Ndrx_bench_write_stats(float64(size),
				float64(requests_to_do)/float64(w.GetDetlaSec())); SUCCEED != ret {
				ctx.TpLogError("Failed to write benchmark stats!")
				return FAIL
			}
		}

	}

	return SUCCEED
}
