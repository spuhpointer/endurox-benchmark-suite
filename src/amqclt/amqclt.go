package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	stomp "github.com/go-stomp/stomp"

	b "exbench"

	atmi "github.com/endurox-dev/endurox-go"
)

var M_ctx *atmi.ATMICtx

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func runbench() int {

	conn, err := stomp.Dial("tcp", "localhost:61613")

	defer conn.Disconnect()

	failOnError(err, "Failed to connect to ActiveMQ")

	sub, err := conn.Subscribe("/queue/cltrply", stomp.AckAuto)

	failOnError(err, "Failed to connect to subscribe to reply queue")

	ret := b.Ndrx_bench_clmain(M_ctx, 1, func(ctx *atmi.ATMICtx, correl int64,
		buf []byte, oneway bool) (int, []byte) {
		corrId := strconv.FormatInt(correl, 10)

		if oneway {

			err := conn.Send(
				"/queue/srvreq", // destination
				"text/plain",    // content-type
				buf,             // body
				stomp.SendOpt.Header("corrid", corrId))

			failOnError(err, "Failed to publish a message")

			return atmi.SUCCEED, nil
		} else {

			err := conn.Send(
				"/queue/srvreq", // destination
				"text/plain",    // content-type
				buf,             // body
				stomp.SendOpt.Header("corrid", corrId))

			failOnError(err, "Failed to publish a message")

			for {
				msg := <-sub.C
				if msg.Header.Get("corrid") == corrId {
					return atmi.SUCCEED, msg.Body
				}
			}

			/* check the header */
			return atmi.FAIL, nil
		}
	})

	return ret
}

func main() {

	var err atmi.ATMIError
	M_ctx, err = atmi.NewATMICtx()

	if nil != err {
		fmt.Fprintf(os.Stderr, "TESTERROR ! Failed to allocate cotnext %s!\n", err)
		os.Exit(atmi.FAIL)
	}

	ret := runbench()

	M_ctx.TpLogInfo("Benchmark finished with %d", ret)

	os.Exit(ret)

}
