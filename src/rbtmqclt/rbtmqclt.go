package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/streadway/amqp"

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
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

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

	/*
	 * So resources are open... now run the code..
	 */

	/*
		corrId := randomString(32)

		err = ch.Publish(
			"",          // exchange
			"rpc_queue", // routing key
			false,       // mandatory
			false,       // immediate
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: corrId,
				ReplyTo:       q.Name,
				Body:          []byte(strconv.Itoa(n)),
			})
		failOnError(err, "Failed to publish a message")

		for d := range msgs {
			if corrId == d.CorrelationId {
				res, err = strconv.Atoi(string(d.Body))
				failOnError(err, "Failed to convert body to integer")
				break
			}
		}
	*/

	ret := b.Ndrx_bench_clmain(M_ctx, 1, func(ctx *atmi.ATMICtx, correl int64, buf []byte) (int, []byte) {
		corrId := strconv.FormatInt(correl, 10)

		err = ch.Publish(
			"",          // exchange
			"rpc_queue", // routing key
			false,       // mandatory
			false,       // immediate
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: corrId,
				ReplyTo:       q.Name,
				Body:          buf,
			})

		failOnError(err, "Failed to publish a message")

		for d := range msgs {
			if corrId == d.CorrelationId {
				//res, err = strconv.Atoi(string(d.Body))
				//failOnError(err, "Failed to convert body to integer")
				return atmi.SUCCEED, d.Body
				break
			}
		}

		ctx.TpLogError("Failed to get response!")
		return atmi.FAIL, nil
	})

	return ret
}

func main() {

	/*
		n := bodyFrom(os.Args)

		log.Printf(" [x] Requesting fib(%d)", n)
		res, err := fibonacciRPC(n)
		failOnError(err, "Failed to handle RPC request")

		log.Printf(" [.] Got %d", res)
	*/

	var err atmi.ATMIError
	M_ctx, err = atmi.NewATMICtx()

	if nil != err {
		fmt.Fprintf(os.Stderr, "TESTERROR ! Failed to allocate cotnext %s!\n", err)
		os.Exit(atmi.FAIL)
	}

	//ret := b.Ndrx_bench_clmain(M_ctx, 1, request)

	ret := runbench()

	M_ctx.TpLogInfo("Benchmark finished with %d", ret)

	os.Exit(ret)

}

func bodyFrom(args []string) int {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "30"
	} else {
		s = strings.Join(args[1:], " ")
	}
	n, err := strconv.Atoi(s)
	failOnError(err, "Failed to convert arg to integer")
	return n
}
