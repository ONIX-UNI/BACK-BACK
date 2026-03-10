package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/DuvanRozoParra/sicou/pkg/messaging"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer := messaging.NewKafkaConsumer(
		[]string{"localhost:9092"},
		"hello-topic",
		"hello-group",
	)
	// "hello-group",
	defer consumer.Close()

	if err := consumer.HealthCheck(ctx); err != nil {
		log.Fatalf("Kafka no disponible: %v", err)
	}

	log.Println("Kafka conectado correctamente")

	go func() {
		err := consumer.Consume(ctx, func(data []byte) error {
			log.Printf("Evento recibido: %s\n", string(data))
			return nil
		})
		if err != nil {
			log.Println("consumer stopped:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	cancel()
}
