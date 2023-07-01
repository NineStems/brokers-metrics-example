package domain

import "time"

type RabbitMessage struct {
	Id   int       `json:"id"`
	Time time.Time `json:"time"`
}

type KafkaMessage struct {
	Topic   string `json:"topic"`
	Content []byte `json:"content"`
}
