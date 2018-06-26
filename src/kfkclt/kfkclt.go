/*
 * @brief Kafka client process benchmark
 *
 * @file kfkclt.go
 */
package main

import (
	b "exbench"
	"fmt"
	"os"
	"strconv"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	atmi "github.com/endurox-dev/endurox-go"
)

var M_ctx *atmi.ATMICtx

var M_consumer *kafka.Consumer
var M_producer *kafka.Producer

var M_request_topic = "srvreq"
var M_reply_topic = "cltrply"

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

	corrId := strconv.FormatInt(correl, 10)

	if oneway {

		//ctx.TpLogInfo("About to produce...")

		if err := M_producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &M_request_topic, Partition: kafka.PartitionAny},
			Value:          buf,
			Key:            []byte(corrId),
		}, nil); nil != err {
			ctx.TpLogError("Failed to produce message: %s", err.Error())
			return atmi.FAIL, nil
		}

		/* we are ok, buffer receive, lets return it... */
		return atmi.SUCCEED, nil

	} else {

		if err := M_producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &M_request_topic, Partition: kafka.PartitionAny},
			Value:          buf,
			Key:            []byte(corrId),
		}, nil); nil != err {
			ctx.TpLogError("Failed to produce message: %s", err.Error())
			return atmi.FAIL, nil
		}

		/* wait for correlated message... */
		for {
			msg, err := M_consumer.ReadMessage(-1)
			if err == nil {
				/*
					ctx.TpLogInfo("Message on %s: %s=%s\n", msg.TopicPartition,
						string(msg.Key), string(msg.Value))
				*/
				if string(msg.Key) == corrId {

					return atmi.SUCCEED, msg.Value
				}
			} else {
				ctx.TpLogError("Consumer error: %v (%v)", err, msg)
				return atmi.FAIL, nil
			}
		}
	}

}

func main() {

	var errA atmi.ATMIError
	var err error

	M_ctx, errA = atmi.NewATMICtx()

	if nil != errA {
		fmt.Fprintf(os.Stderr, "TESTERROR ! Failed to allocate cotnext %s!\n", errA)
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

	M_consumer.SubscribeTopics([]string{M_reply_topic}, nil)

	//////////////////////////////////////////////////////////////////////////////
	// Create producer
	//////////////////////////////////////////////////////////////////////////////
	M_producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		panic(err)
	}

	ret := b.Ndrx_bench_clmain(M_ctx, 1, request)

	M_ctx.TpLogInfo("Benchmark finished with %d", ret)

	M_consumer.Close()
	M_producer.Close()

	os.Exit(ret)
}
