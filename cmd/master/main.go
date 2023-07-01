package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"net/http"
)

var (
	rabbitProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "rabbit_processed",
		Help: "The total number of rabbit messages",
	})
	kafkaProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "kafka_processed",
		Help: "The total number of kafka messages",
	})
)

func rabbit(w http.ResponseWriter, req *http.Request) {
	defer w.WriteHeader(200)
	body, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("error by read rabbit %v", err.Error())
		return
	}
	defer req.Body.Close()
	fmt.Printf("message from rabbit %v\n", string(body))
	rabbitProcessed.Inc()
}

func kafka(w http.ResponseWriter, req *http.Request) {
	defer w.WriteHeader(200)
	body, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("error by read kafka %v", err.Error())
		return
	}
	defer req.Body.Close()
	fmt.Printf("message from kafka %v\n", string(body))
	kafkaProcessed.Inc()
}

func main() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/rabbit", rabbit)
	http.HandleFunc("/kafka", kafka)
	http.ListenAndServe("localhost:8000", nil)
}
