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

	first := true
	conn, err := stomp.Dial("tcp", "localhost:61613")

	defer conn.Disconnect()

	failOnError(err, "Failed to connect to ActiveMQ")

	sub, err := conn.Subscribe("/queue/cltrply", stomp.AckAuto)

	failOnError(err, "Failed to connect to subscribe to reply queue")

	ret := b.Ndrx_bench_clmain(M_ctx, 1, func(ctx *atmi.ATMICtx, correl int64,
		buf []byte, oneway bool) (int, []byte) {
		corrId := strconv.FormatInt(correl, 10)

		if oneway {

			if first {
				sub.Unsubscribe()
			}

			//We get here "2018/06/27 09:03:04 Failed to publish a message: connection already closed"
			//Try to reconnect on error. Quick and direty solution...
		restart:
			err := conn.Send(
				"/queue/srvreq",            // destination
				"application/octet-stream", // content-type
				buf, // body
				stomp.SendOpt.Header("corrid", corrId))

			if nil != err {
				M_ctx.TpLogError("Got error: %s", err.Error())
				//Try to reconnect..
				conn, err = stomp.Dial("tcp", "localhost:61613")

				failOnError(err, "Failed to connect to ActiveMQ")

				goto restart
			}

			failOnError(err, "Failed to publish a message")

			return atmi.SUCCEED, nil
		} else {

			//M_ctx.TpLogInfo("Sending corr [%s] %s", corrId, string(buf))
			err := conn.Send(
				"/queue/srvreq", // destination
				"text/plain",    // content-type
				buf,             // body
				stomp.SendOpt.Header("corrid", corrId))

			failOnError(err, "Failed to publish a message")

			//M_ctx.TpLogInfo("Waiting for reply...")

			for {
				msg := <-sub.C

				//M_ctx.TpLogInfo("Got corr [%s] %s", msg.Header.Get("corrid"), string(buf))

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
