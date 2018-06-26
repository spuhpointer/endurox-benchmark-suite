package main

import (
	b "exbench"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	atmi "github.com/endurox-dev/endurox-go"
)

var M_ctx *atmi.ATMICtx

var M_consumer *kafka.Consumer

var M_request_topic = "srvreq"

func main() {

	var err error
	//Have some context
	M_ctx, errA := atmi.NewATMICtx()

	if nil != errA {
		fmt.Fprintf(os.Stderr, "Failed to allocate new context: %s", errA)
		os.Exit(atmi.FAIL)
	}

	//////////////////////////////////////////////////////////////////////////////
	//create consumer
	//////////////////////////////////////////////////////////////////////////////
	M_consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	M_consumer.SubscribeTopics([]string{M_request_topic}, nil)

	defer M_consumer.Close()

	/* wait for correlated message... */
	for {
		//M_ctx.TpLogInfo("Waiting for msgs...")
		msg, err := M_consumer.ReadMessage(-1)

		if err == nil {

			//M_ctx.TpLogInfo("got msg: %s", string(msg.Value))
			ret := b.Ndrx_bench_svmain_oneway(M_ctx, 0, msg.Value)

			if ret != atmi.SUCCEED {
				M_ctx.TpLogError("Failed to process incoming message!")
				os.Exit(atmi.FAIL)
			}

		} else {
			M_ctx.TpLogError("Consumer error: %v (%v)", err, msg)
			os.Exit(atmi.FAIL)
		}
	}

	M_ctx.TpLogInfo("Benchmark finished")

	os.Exit(atmi.SUCCEED)

}
