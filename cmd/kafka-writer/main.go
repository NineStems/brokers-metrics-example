package main

import (
	"context"
	"fmt"
	"mb-and-metrics/internal/config"
	"mb-and-metrics/internal/pkg/logger"
	"mb-and-metrics/internal/repositories/kafka"
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

	producer, err := kafka.NewProducer(strings.Split(cfg.Kafka.Brokers, ","))
	if err != nil {
		sugarLog.Error(err)
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Minute*1)

	var counter int
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Second * 1)
		}
		counter++
		if err = producer.Send(cfg.Kafka.Topic, fmt.Sprintf("id:%d time:%v", counter, time.Now())); err != nil {
			sugarLog.Error(err)
		}
	}

}
