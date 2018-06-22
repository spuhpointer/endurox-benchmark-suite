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

import atmi "github.com/endurox-dev/endurox-go"

/*
import "unsafe"
import "fmt"
import "runtime"
*/

/**
 * SUCCEED/FAIL flags
 */
const (
	/** succeed status */
	SUCCEED = 0
	/** fail status */
	FAIL = -1
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
 * @return -1 all OK, or >=0, position at which value is invalid + (expected, got)
 */
func Ndrx_bench_verify_buffer(buf []byte, size int, mod int) (int, byte, byte) {

	for i := 0; i < size; i++ {
		expected := byte(255 - i%mod)

		if expected != buf[i] {
			return i, byte(expected), buf[i]
		}
	}
	return -1, 0, 0
}

/**
 * request callback function
 * @param ctx[in] ATMI context (for logging, etc...)
 * @param buf[in] Byte buffer to send
 * @param correl[in] Call correlator (used by systems where needed to match req with rsp)
 * @return status code 0 = succeed, -1 = FAIL, return buffer.
 */
type Ndrx_Bench_requestCB func(ctx atmi.ATMICtx, correl int64, buf []byte) (int, []byte)

/**
 * Benchmark main
 * @param threadid[in] logical thread id used for correlators...
 * @param request[in] callback function for sending the data to server and
 *   receiving response back
 * @return 0 = succeed, -1 fail
 */
func Ndrx_bench_main(threadid int, nrrequests int, request Ndrx_Bench_requestCB) int {

	ctx, err := atmi.NewATMICtx()

	if nil != err {
		fmt.Errorf("Failed to allocate cotnext!", err)
		return FAIL
	}

	size := 0

	/* Correlator shall be built as thread id + call number (we need some offset
	 * for negative steps)
	 */
	for i := -1; i < 56; i++ {

		if i <= 0 {
			size = 1
		} else {
			size = i * 1024
		}

		ctx.TpLogInfo("Benchmarking step: %d with size of %d bytes (nr req: %d)",
			i, size, nrrequests)

		/* we shall loop over the given count, let say nrrequests / 100K requests
		 * and shall start the stopwatch
		 */

		for req := 0; req < nrrequests; req++ {

			buf := Ndrx_bench_get_buffer(size, 255)

			/* build up correlator... */
			var correl int64

			correl = threadid

			correl <<= 8

			correl |= i

			correl <<= 40

			correl |= req

			res, retbuf := request(ctx, correl, buf)
		}

	}

	return SUCCEED
}
