package dao

import (
	"log"
	"replite_web/internal/app/config"
	"strings"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
)

var singleConnsPool map[string]*kafka.Writer

func newKafkaWriter(topic string) *kafka.Writer {
	topics := config.ServerConf.Kafka.Topic
	var find = false
	for i := 0; i < len(topics); i++ {
		tmp := topics[i]
		if strings.Compare(tmp, topic) == 0 {
			find = true
			break
		}
	}
	if !find {
		log.Println("failed use the topic which cant configurate in the server.xml")
		return nil
	}
	config := kafka.WriterConfig{
		Brokers: config.ServerConf.Kafka.Broker,
		Topic:   topic,
	}
	writer := kafka.NewWriter(config)
	writer.AllowAutoTopicCreation = true
	writer.Compression = compress.Snappy
	return writer
}

func GetTopicConn(topic string) *kafka.Writer {
	if writer, ok := singleConnsPool[topic]; ok {
		return writer
	}
	writer := newKafkaWriter(topic)
	singleConnsPool[topic] = writer
	return writer
}
