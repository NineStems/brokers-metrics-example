package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"mb-and-metrics/domain"
	"sync"
)

var brokers = []string{"127.0.0.1:9092"}

type Consumer struct {
	consumer sarama.Consumer
}

func NewConsumer(brokers []string) (*Consumer, error) {
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, err
	}

	return &Consumer{consumer: consumer}, nil
}

func (c *Consumer) Subscribe(topic string) (chan domain.KafkaMessage, error) {
	partitionList, err := c.consumer.Partitions(topic) //get all partitions on the given topic
	if err != nil {
		return nil, fmt.Errorf("retrieving partitionList: %w", err)
	}
	initialOffset := sarama.OffsetOldest //get offset for the oldest message on t
	out := make(chan domain.KafkaMessage)
	merge := make(chan *sarama.ConsumerMessage)
	go func() {
		defer close(out)
		for message := range merge {
			out <- domain.KafkaMessage{
				Topic:   message.Topic,
				Content: message.Value,
			}
		}
	}()

	go func() {
		defer close(merge)
		wg := sync.WaitGroup{}
		for _, partition := range partitionList {
			pc, err := c.consumer.ConsumePartition(topic, partition, initialOffset)
			if err != nil {
				log.Printf("ConsumePartition: %w", err)
				return
			}

			wg.Add(1)
			go func(pc sarama.PartitionConsumer, wg *sync.WaitGroup) {
				defer wg.Done()
				for message := range pc.Messages() {
					merge <- message
				}
			}(pc, &wg)
		}
		wg.Wait()
	}()

	return out, nil
}

type Producer struct {
	producer sarama.SyncProducer
}

func NewProducer(brokers []string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config)

	return &Producer{producer: producer}, err
}

func (p *Producer) Send(topic string, message string) error {
	partition, offset, err := p.producer.SendMessage(prepareMessage(topic, message))
	if err != nil {
		return fmt.Errorf("%s error occured", err.Error())
	}
	log.Printf("Message was saved to partion: %d.\nMessage offset is: %d.\n", partition, offset)
	return nil
}

func prepareMessage(topic, message string) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	return msg
}
