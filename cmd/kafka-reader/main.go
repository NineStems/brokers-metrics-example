package main

import (
	"bytes"
	"context"
	"mb-and-metrics/internal/config"
	"mb-and-metrics/internal/pkg/logger"
	"mb-and-metrics/internal/repositories/kafka"
	"net/http"
	"strings"
	"time"
)

func main() {
	cfg := config.New()
	err := cfg.Apply("configs/config.yaml")
	if err != nil {
		panic(err)
	}

	log := logger.Console(cfg.Logger.Path, cfg.Logger.Level)
	sugarLog := logger.InitSugarZapLogger(log)

	consumer, err := kafka.NewConsumer(strings.Split(cfg.Kafka.Brokers, ","))
	ctx, _ := context.WithTimeout(context.Background(), time.Minute*1)

	messages, err := consumer.Subscribe(cfg.Kafka.Topic)
	if err != nil {
		sugarLog.Error(err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case message := <-messages:
			sugarLog.Info("read from kafka:", string(message.Content))
			sendReq(cfg.Kafka.Main, message.Content, sugarLog)
		default:
			time.Sleep(time.Millisecond * 1000)
		}
	}
}

func sendReq(address string, body []byte, sugarLog *logger.Logger) {
	req, err := http.NewRequest(http.MethodPost, address, bytes.NewReader(body))
	if err != nil {
		sugarLog.Error(err)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		sugarLog.Error(err)
		return
	}

	defer res.Body.Close()
}
