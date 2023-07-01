package main

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"mb-and-metrics/domain"
	"mb-and-metrics/internal/config"
	"mb-and-metrics/internal/pkg/logger"
	"mb-and-metrics/internal/repositories/rabbit"
	"net"
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

	connRabbit, err := amqp.Dial(fmt.Sprintf(
		"amqp://%v:%v@%v/",
		cfg.Rabbit.Credential.Username,
		cfg.Rabbit.Credential.Password,
		net.JoinHostPort(cfg.Rabbit.Host, cfg.Rabbit.Port),
	))
	if err != nil {
		sugarLog.Error(err)
		return
	}

	defer connRabbit.Close()

	rabbitClient := rabbit.New(connRabbit, sugarLog, cfg.Rabbit)
	ctx, _ := context.WithTimeout(context.Background(), time.Minute*1)

	if err = rabbitClient.Start(); err != nil {
		sugarLog.Error(err)
	}

	var counter int
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Second * 3)
		}
		counter++
		message := domain.RabbitMessage{
			Id:   counter,
			Time: time.Now(),
		}

		if err = rabbitClient.Publish(ctx, message); err != nil {
			sugarLog.Error(err)
		}
	}

}
