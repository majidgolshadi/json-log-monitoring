package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Shopify/sarama"
	"github.com/wvanbergen/kafka/consumergroup"
	"github.com/wvanbergen/kazoo-go"
	"github.com/majidgolshadi/json-log-monitoring"
)

type Config struct {
	Kafka_consumer KafkaConsumerConfig
	Monitoring Monitoring
	Port string
}

type KafkaConsumerConfig struct {
	Zookeeper string
	Topics    string
	GroupName string `toml:"consumer_group"`
	CommitBuffer    int `toml:"commit_buffer"`
}

type Monitoring struct {
	CountingRegex string `toml:"counting_regex"`
}

func main() {
	if _, err := os.Stat("config.toml"); err != nil {
		log.Fatal("Config file is missing")
		os.Exit(2)
	}

	cnf := &Config{}
	if _, err := toml.DecodeFile("config.toml", cnf); err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

	logger := log.New(os.Stdout, "[kafka-consumer] ", log.LstdFlags)
	sarama.Logger = logger

	consumer, err := initConsumer(cnf.Kafka_consumer)
	if err != nil {
		logger.Fatal(err.Error())
		return
	}

	defer consumer.Close()
	setupInterruptListener(consumer)

	analyzer := json_log_monitoring.CreateAnalyzer(cnf.Monitoring.CountingRegex)

	go func() {
		cb :=  cnf.Kafka_consumer.CommitBuffer
		for message := range consumer.Messages() {
			analyzer.Analyze(message.Value)
			cb--
			if cb < 1 {
				consumer.CommitUpto(message)
				cb = cnf.Kafka_consumer.CommitBuffer
			}
		}
	}()

	json_log_monitoring.RunHttpServer(analyzer, cnf.Port)
}

func setupInterruptListener(consumer *consumergroup.ConsumerGroup) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		if err := consumer.Close(); err != nil {
			println("Error closing the consumer", err.Error())
		}
	}()
}

func initConsumer(kc KafkaConsumerConfig) (*consumergroup.ConsumerGroup, error) {
	var zookeeperNodes []string

	config := consumergroup.NewConfig()
	config.Offsets.Initial = sarama.OffsetNewest
	config.Offsets.ProcessingTimeout = 10 * time.Second
	zookeeperNodes, config.Zookeeper.Chroot = kazoo.ParseConnectionString(kc.Zookeeper)
	kafkaTopics := strings.Split(kc.Topics, ",")

	return consumergroup.JoinConsumerGroup(kc.GroupName, kafkaTopics, zookeeperNodes, config)
}
