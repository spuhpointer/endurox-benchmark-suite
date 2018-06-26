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

	//Have some context
	M_ctx, errA := atmi.NewATMICtx()

	if nil != errA {
		fmt.Fprintf(os.Stderr, "Failed to allocate new context: %s", errA)
		os.Exit(atmi.FAIL)
	}

	conn, err := stomp.Dial("tcp", "localhost:61613")

	defer conn.Disconnect()

	failOnError(err, "Failed to connect to ActiveMQ")

	sub, err := conn.Subscribe("/queue/srvreq", stomp.AckAuto)

	failOnError(err, "Failed to connect to subscribe to reply queue")

	for {
		msg := <-sub.C

		ret, ret_bytes := b.Ndrx_bench_svmain(M_ctx, 0, msg.Body)

		if ret != atmi.SUCCEED {
			M_ctx.TpLogError("Failed to process incoming message!")
			os.Exit(atmi.FAIL)
		}

		err := conn.Send(
			"/queue/cltreply", // destination
			"text/plain",      // content-type
			ret_bytes,         // body
			stomp.SendOpt.Header("corrid", msg.Header.Get("corrid")))

		failOnError(err, "Failed to publish a message")

	}

}
