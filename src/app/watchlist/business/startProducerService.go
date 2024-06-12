package business

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"stock_broker_application/src/utils/postgres"
	"sync"
	"time"
	"watchlist/repositories"

	"github.com/segmentio/kafka-go"
)

func StartProducer() {
	go createComment()
}

func createComment() {
	brokerAddress := os.Getenv("KAFKA_BROKER")
	if brokerAddress == "" {
		log.Fatal("KAFKA_BROKER environment variable is not set")
		return
	}

	topic := "updateStockPrice"
	//rand.Seed(time.Now().UnixNano()) // Seed the random number generator

	// Assume we have a global variable for the database connection
	symbols, err := repositories.NewGetStocksSymbolRepository(postgres.DB).GetStocksSymbols()
	if err != nil {
		log.Println("Error getting stocks symbols from database:", err)
		return
	}

	var wg sync.WaitGroup

	for _, symbol := range symbols {
		wg.Add(1)
		go func(symbol string) {
			defer wg.Done()
			for {
				randomStockPrice := rand.Intn(21) + 90
				err := pushStockPriceToQueue(brokerAddress, topic, symbol, randomStockPrice)
				if err != nil {
					log.Println("Error pushing stock price to queue:", err)
				} else {
					log.Printf("Stock price for %s pushed to queue: %d", symbol, randomStockPrice)
				}
				time.Sleep(10 * time.Second)
			}
		}(symbol)
	}

	wg.Wait()
}

func pushStockPriceToQueue(brokerAddress, topic, symbol string, stockPrice int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{brokerAddress},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	message := kafka.Message{
		Key:   []byte(symbol),
		Value: []byte(fmt.Sprintf("%d", stockPrice)),
	}

	err := writer.WriteMessages(ctx, message)
	if err != nil {
		return fmt.Errorf("could not write message: %w", err)
	}

	log.Printf("Sent: %s: %s\n", symbol, string(message.Value))
	return nil
}
