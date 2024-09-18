package main

import (
	"crypto/md5"
	"encoding/hex"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"md5_queue", gitit
	true,
		false,
		false,
		false,
		nil,
)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	log.Println("Waiting for messages. To exit press CTRL+C")
	for msg := range msgs {
		hash := GenerateMD5(msg.Body)
		log.Printf("Processed message: %s | MD5 Hash: %s\n", msg.Body, hash)
	}
}

func GenerateMD5(input []byte) string {
	hash := md5.Sum(input)
	return hex.EncodeToString(hash[:])
}
