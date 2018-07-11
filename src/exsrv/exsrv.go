/**
 * @brief Enduro/X server process
 *
 * @file exsrv.go
 */
package main

import (
	b "exbench"
	"fmt"
	"os"

	atmi "github.com/endurox-dev/endurox-go"
)

const (
	UNSET   = -1
	FALSE   = 0
	TRUE    = 1
	SUCCEED = atmi.SUCCEED
	FAIL    = atmi.FAIL
)

//Un-init function
func appinit(ctx *atmi.ATMICtx) int {

	if errA := ctx.TpAdvertise("EXSVC", "EXSVC", EXSVC); nil != errA {

		ctx.TpLogError("Failed to advertise EXSVC: %s", errA.Message())
		return FAIL
	}

	if errA := ctx.TpAdvertise("EXONEWAY", "EXONEWAY", EXONEWAY); nil != errA {

		ctx.TpLogError("Failed to advertise EXONEWAY: %s", errA.Message())
		return FAIL
	}

	return SUCCEED
}

//One way service (verify incomding data and perform timings...
//@param ac ATMI Context
//@param svc Service call information
func EXONEWAY(ac *atmi.ATMICtx, svc *atmi.TPSVCINFO) {

	ret := SUCCEED

	buf, errA := ac.CastToCarray(&svc.Data)

	if nil != errA {
		ac.TpLogError("Failed to cast to Carray: %s", errA.Message())
		ac.TpReturn(atmi.TPFAIL, 0, svc.Data.GetBuf(), 0)
		return
	}

	//Run off the bencharmk suite
	ret = b.Ndrx_bench_svmain_oneway(ac, 0, buf.GetBytes())

	if ret != SUCCEED {
		ac.TpLogError("Failed to process incoming message! ")
		ac.TpReturn(atmi.TPFAIL, 0, svc.Data.GetBuf(), 0)
		return
	}

	/* does not send data back anyway... */
	ac.TpReturn(atmi.TPSUCCESS, 0, buf, 0)

	return
}

//EXSVC service - generic entry point
//@param ac ATMI Context
//@param svc Service call information
func EXSVC(ac *atmi.ATMICtx, svc *atmi.TPSVCINFO) {

	ret := SUCCEED

	buf, errA := ac.CastToCarray(&svc.Data)

	if nil != errA {
		ac.TpLogError("Failed to cast to Carray: %s", errA.Message())
		ac.TpReturn(atmi.TPFAIL, 0, svc.Data.GetBuf(), 0)
		return
	}

	//Run off the bencharmk suite
	ret, ret_bytes := b.Ndrx_bench_svmain(ac, 0, buf.GetBytes())

	if ret != SUCCEED {
		ac.TpLogError("Failed to process incoming message! ")
		ac.TpReturn(atmi.TPFAIL, 0, svc.Data.GetBuf(), 0)
		return
	}

	buf.SetBytes(ret_bytes)

	ac.TpReturn(atmi.TPSUCCESS, 0, buf, 0)

	return
}

//Un-init & Terminate the application
func unInit(ac *atmi.ATMICtx) {

	ac.TpLogInfo("Shutdown ok")
}

//Executable main entry point
func main() {
	//Have some context
	ac, err := atmi.NewATMICtx()

	if nil != err {
		fmt.Fprintf(os.Stderr, "Failed to allocate new context: %s", err)
		os.Exit(atmi.FAIL)
	} else {
		//Run as server
		if err = ac.TpRun(appinit, unInit); nil != err {
			ac.TpLogError("Exit with failure")
			os.Exit(atmi.FAIL)
		} else {
			ac.TpLogInfo("Exit with success")
			os.Exit(atmi.SUCCEED)
		}
	}
}
