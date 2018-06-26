
# Create kafka topics

```
sudo /opt/Kafka/kafka_2.10-0.10.0.1/bin/kafka-topics.sh --create --zookeeper localhost:2181 --replication-factor 1  --partitions 1 --topic cltrply
sudo /opt/Kafka/kafka_2.10-0.10.0.1/bin/kafka-topics.sh --create --zookeeper localhost:2181 --replication-factor 1  --partitions 1 --topic srvreq
sudo /opt/Kafka/kafka_2.10-0.10.0.1/bin/kafka-topics.sh --zookeeper localhost:2181 --alter --topic cltrply --config retention.ms=1000
sudo /opt/Kafka/kafka_2.10-0.10.0.1/bin/kafka-topics.sh --zookeeper localhost:2181 --alter --topic srvreq --config retention.ms=1000
```

