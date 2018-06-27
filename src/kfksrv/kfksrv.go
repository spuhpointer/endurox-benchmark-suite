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
var M_producer *kafka.Producer

var M_request_topic = "cltrply"
var M_reply_topic = "srvreq"

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

	M_consumer.SubscribeTopics([]string{M_reply_topic}, nil)

	defer M_consumer.Close()

	//////////////////////////////////////////////////////////////////////////////
	// Create producer
	//////////////////////////////////////////////////////////////////////////////
	M_producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		panic(err)
	}

	defer M_producer.Close()

	/* wait for correlated message... */
	for {
		msg, err := M_consumer.ReadMessage(-1)
		if err == nil {

			//M_ctx.TpLogInfo("Message on %s: %s=%s\n", msg.TopicPartition,
			//	string(msg.Key), string(msg.Value))
			ret, ret_bytes := b.Ndrx_bench_svmain(M_ctx, 0, msg.Value)

			if ret != atmi.SUCCEED {
				M_ctx.TpLogError("Failed to process incoming message!")
				os.Exit(atmi.FAIL)
			}

			if err := M_producer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &M_request_topic, Partition: kafka.PartitionAny},
				Value:          ret_bytes,
				Key:            msg.Key,
			}, nil); nil != err {
				M_ctx.TpLogError("Failed to produce message: %s", err.Error())
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
