# Do recursive builds
all:
	go get -u github.com/endurox-dev/endurox-go
	go get -u github.com/streadway/amqp
	$(MAKE) -C exbench
	$(MAKE) -C exsrv
	$(MAKE) -C exclt
	$(MAKE) -C rbtmqsrv
	$(MAKE) -C rbtmqsrvoneway
	$(MAKE) -C rbtmqclt
	$(MAKE) -C kfkclt
	$(MAKE) -C kfksrv

clean:
	$(MAKE) -C exbench clean
	$(MAKE) -C exsrv clean
	$(MAKE) -C exclt clean
	$(MAKE) -C rbtmqsrv clean
	$(MAKE) -C rbtmqsrvoneway clean
	$(MAKE) -C rbtmqclt clean
	$(MAKE) -C kfkclt clean
	$(MAKE) -C kfksrv clean

.PHONY: clean
