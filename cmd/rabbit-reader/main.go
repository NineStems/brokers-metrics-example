package main

import (
	"bytes"
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"mb-and-metrics/internal/config"
	"mb-and-metrics/internal/pkg/logger"
	"mb-and-metrics/internal/repositories/rabbit"
	"net"
	"net/http"
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

	var message = make(chan []byte)
	defer close(message)
	if err = rabbitClient.Read(ctx, message); err != nil {
		sugarLog.Error(err)
	}

	var info []byte
	for {
		info = nil
		select {
		case <-ctx.Done():
			return
		case info = <-message:
			sugarLog.Info("read from rabbit:", string(info))
			sendReq(cfg.Rabbit.Main, info, sugarLog)
		default:
			time.Sleep(time.Millisecond * 100)
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
