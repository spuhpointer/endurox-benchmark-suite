/*
 * @brief Enduro/X Middleware client process benchmark
 *
 * @file exclt.go
 */
package main

import (
	b "exbench"
	"fmt"
	"os"

	atmi "github.com/endurox-dev/endurox-go"
)

var M_ctx *atmi.ATMICtx

const TESTSVC = "EXSVC"
const TESTSVC1W = "EXONEWAY"

/**
 * Send message to server process
 * @param ctx ATMI Context
 * @param correl correlator
 * @param buf Buffer to send
 * @return status code
 * @return if status ok, byte array in response
 */
func request(ctx *atmi.ATMICtx, correl int64, buf []byte, oneway bool) (int, []byte) {

	//Maybe re-use the buffer?
	carray, errA := ctx.NewCarray(buf)

	if nil != errA {
		ctx.TpLogError("Failed to allocate new Carray bufer: %s", errA.Message())
		return atmi.FAIL, nil
	}

	if oneway {

		_, errA = ctx.TpACall(TESTSVC1W, carray, atmi.TPNOREPLY)

		if nil != errA {
			ctx.TpLogError("Failed to call [%s] service: %s", TESTSVC1W, errA.Message())
			return atmi.FAIL, nil
		}

		/* we are ok, buffer receive, lets return it... */
		return atmi.SUCCEED, nil

	} else {

		_, errA = ctx.TpCall(TESTSVC, carray, 0)

		if nil != errA {
			ctx.TpLogError("Failed to call [%s] service: %s", TESTSVC, errA.Message())
			return atmi.FAIL, nil
		}

		/* we are ok, buffer receive, lets return it... */

		return atmi.SUCCEED, carray.GetBytes()
	}

}

func main() {

	var err atmi.ATMIError
	M_ctx, err = atmi.NewATMICtx()

	if nil != err {
		fmt.Fprintf(os.Stderr, "TESTERROR ! Failed to allocate cotnext %s!\n", err)
		os.Exit(atmi.FAIL)
	}

	ret := b.Ndrx_bench_clmain(M_ctx, 1, request)

	M_ctx.TpLogInfo("Benchmark finished with %d", ret)

	os.Exit(ret)
}
