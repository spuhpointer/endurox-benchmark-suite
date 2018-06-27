package main

import (
	b "exbench"
	"fmt"
	"os"

	atmi "github.com/endurox-dev/endurox-go"
	"github.com/streadway/amqp"
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

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"rpc_queue", // name
		false,       // durable
		false,       // delete when usused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {

			//Run off the bencharmk suite
			ret, ret_bytes := b.Ndrx_bench_svmain(M_ctx, 0, d.Body)

			if ret != atmi.SUCCEED {
				M_ctx.TpLogError("Failed to process incoming message!")
				os.Exit(atmi.FAIL)
			}

			err = ch.Publish(
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "application/octet-stream",
					CorrelationId: d.CorrelationId,
					Body:          ret_bytes,
				})
			failOnError(err, "Failed to publish a message")

			// Auto ack will be ok, because we return RPC reply, thus ack
			//d.Ack(false)
		}
	}()

	<-M_quit
}
