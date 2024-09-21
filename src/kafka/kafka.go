package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
)

var (
	config *sarama.Config
	producer sarama.AsyncProducer
)

func init() {
	config = sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	var err error
	producer, err = sarama.NewAsyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		panic(err)
	}
}

func WriteToKafka(msg string) {

	defer func() {
		if err := producer.Close(); err != nil {
			panic(err)
		}
	}()

	message := &sarama.ProducerMessage{Topic: "test", Value: sarama.StringEncoder(msg)}

	producer.Input() <- message

	select {
	case success := <-producer.Successes():
		fmt.Println("Message produced:", success.Offset)
	case err := <-producer.Errors():
		fmt.Println("Failed to produce message:", err)
	}
}
