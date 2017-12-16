package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
	"encoding/json"

	"github.com/BurntSushi/toml"
	"github.com/Shopify/sarama"
	"github.com/wvanbergen/kafka/consumergroup"
	"github.com/wvanbergen/kazoo-go"
	"github.com/majidgolshadi/json-log-monitoring/rest"
)

type Config struct {
	Kafka_consumer KafkaConsumerConfig
	Port string
}

type KafkaConsumerConfig struct {
	Zookeeper string
	Topics    string
	GroupName string `toml:"customer_group"`
	CommitBuffer    int `toml:"commit_buffer"`
}

var applicationConfig *Config

func init() {
	if _, err := os.Stat("config.toml"); err != nil {
		log.Fatal("Config file is missing")
		os.Exit(2)
	}

	applicationConfig = &Config{}
	if _, err := toml.DecodeFile("config.toml", applicationConfig); err != nil {
		log.Fatal(err)
		os.Exit(2)
	}
}

func main() {
	logger := log.New(os.Stdout, "[kafka-consumer] ", log.LstdFlags)
	sarama.Logger = logger

	consumer, err := initConsumer(applicationConfig.Kafka_consumer)

	defer consumer.Close()

	if err != nil {
		logger.Fatal(err.Error())
		return
	}

	setupInterruptListener(consumer)

	go rest.RunHttpServer(applicationConfig.Port)

	var jsonStruct map[string]interface{}
	cb :=  applicationConfig.Kafka_consumer.CommitBuffer

	for message := range consumer.Messages() {
		if json.Unmarshal(message.Value, &jsonStruct) != nil {
			rest.JsonStructErrorflag = true
		}

		cb--

		if cb < 1 {
			consumer.CommitUpto(message)
			cb = applicationConfig.Kafka_consumer.CommitBuffer
		}
	}
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
