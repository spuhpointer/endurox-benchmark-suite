package main

import (
	b "exbench"
	"fmt"
	"os"

	atmi "github.com/endurox-dev/endurox-go"
	stomp "github.com/go-stomp/stomp"
)

var M_ctx *atmi.ATMICtx

var M_quit = make(chan struct{})

func failOnError(err error, msg string) {
	if err != nil {
		M_ctx.TpLogError("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {

	var errA atmi.ATMIError

	//Have some context
	M_ctx, errA = atmi.NewATMICtx()

	if nil != errA {
		fmt.Fprintf(os.Stderr, "Failed to allocate new context: %s", errA)
		os.Exit(atmi.FAIL)
	}

	if nil != errA {
		fmt.Fprintf(os.Stderr, "Failed to allocate new context: %s", errA)
		os.Exit(atmi.FAIL)
	}

	conn, err := stomp.Dial("tcp", "localhost:61613")

	defer conn.Disconnect()

	failOnError(err, "Failed to connect to ActiveMQ")

	sub, err := conn.Subscribe("/queue/srvreq", stomp.AckAuto)

	failOnError(err, "Failed to connect to subscribe to reply queue")

	M_ctx.TpLogInfo("About waiting for messages...")
	for {
		msg := <-sub.C

		//Run off the bencharmk suite
		ret := b.Ndrx_bench_svmain_oneway(M_ctx, 0, msg.Body)

		if ret != atmi.SUCCEED {
			M_ctx.TpLogError("Failed to process incoming message!")
			os.Exit(atmi.FAIL)
		}

	}
}
