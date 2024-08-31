package events

import (
	"context"
	"log"
	"os"
	"strings"
	cf "structure-golang/config"
	"sync"

	"github.com/IBM/sarama"
)

// Sarama configuration options
var (
	assignor = "range"
	oldest   = true
	verbose  = false
)

// NOTE ADAPTER -------------------------------------
type eventKafka struct {
	brokers  []string
	config   *sarama.Config
	producer sarama.SyncProducer
	// consumer   sarama.ConsumerGroup
}

func NewEventKafka() AppEvent {
	// NOTE Validate Credential -------------------------------------------------
	brokers := strings.Split(cf.Env.KafkaBrokers, ",")
	if len(brokers) == 0 {
		panic("no Kafka bootstrap brokers defined, please set the Brokers in KafkaCredential")
	}

	// NOTE สำหรับ Create Config sarama เพื่อใช้ในการ Connect kafka -----------------
	config := sarama.NewConfig()
	version, err := sarama.ParseKafkaVersion(cf.Env.KafkaVersion)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}
	config.Version = version
	if verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	// NOTE Producer Config ------------------------------------------------------
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewRandomPartitioner // COMMENT Random partition when produce
	producer, err := sarama.NewSyncProducer(brokers, config)  // COMMENT sync producer
	if err != nil {
		panic(err)
	}

	// NOTE Consumer Config -----------------------------------------------------
	switch assignor {
	case "sticky":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	case "roundrobin":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	case "range":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	default:
		log.Panicf("Unrecognized consumer group partition assignor: %s", assignor)
	}

	if oldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	// TODO Return --------------------------------------------------------------
	return &eventKafka{brokers: brokers, config: config, producer: producer}
}

func (e *eventKafka) On(topics string, group string, handler func(topic string, message []byte) error) (err error) {
	// NOTE Validation Credential
	if len(topics) == 0 {
		panic("no topics given to be consumed, please set the topics")
	}

	if len(group) == 0 {
		panic("no Kafka consumer group defined, please set the group")
	}

	keepRunning := true
	log.Printf("Starting a new Sarama consumer topic : '%s'", topics)

	/**
	 * Setup a new Sarama consumer group
	 */
	consumer := Consumer{
		ready:    make(chan bool),
		callback: handler,
	}
	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(e.brokers, group, e.config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	// consumptionIsPaused := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(ctx, strings.Split(topics, ","), &consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}

	}()

	<-consumer.ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	// sigusr1 := make(chan os.Signal, 1)
	// signal.Notify(sigusr1, syscall.S)

	sigterm := make(chan os.Signal, 1)
	// signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			keepRunning = false
			wg.Done()
		}
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}

	return err
}

func (e *eventKafka) Emit(topic string, data []byte) (patition int32, offset int64, err error) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(data),
	}
	patition, offset, err = e.producer.SendMessage(msg)
	if err != nil {
		return 0, 0, err
	}
	return patition, offset, nil
}

// NOTE Helper function -----------------------------------------------------------------------

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready    chan bool
	callback func(topic string, message []byte) error
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message := <-claim.Messages():
			// log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			consumer.callback(message.Topic, message.Value)
			session.MarkMessage(message, "")

		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/Shopify/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}
