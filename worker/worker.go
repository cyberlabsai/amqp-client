package worker

import (
	"fmt"

	"github.com/cyberlabsai/amqp-client/client"
	"github.com/streadway/amqp"
)

// Topics is a type that contains a list of topics and their function.
// The key is the name and the value is what the handler will do when consume a message with this topic.
// Ex:
//   topic: 	"deleted"
//	 function: 	onDelete(msg amqp.Delivery) bool { // delete from storage }
type Topics map[string]func(amqp.Delivery) bool

// EventHandler consume messages and process according their instructions.
// The service help the EventHandler to connect and bind a topic to receive the messages from network.
// The Topics can be added in safe mode with the function AddTopic, or directly.
type EventHandler struct {
	*client.Service
	Topics
}

// New returns a created and connected EventHandler.
// The EventHandler has a service that will declare a new queue
// and link to a topic.
// Before consuming the messages it's necessary to add at least one event
// on your map, after that, the EventHandler can be started successfully.
func New(topicPrefix, address, exchange, contentType, queueName string) (*EventHandler, error) {

	s := client.New(topicPrefix, address, exchange, contentType)
	if err := s.Start(); err != nil {
		return nil, err
	}

	_, err := s.Channel.QueueDeclare(
		queueName,
		false, // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return nil, err
	}

	return &EventHandler{
		Service: s,
	}, nil
}

// Start initiates the EventHandler to consuming and process the messages,
// the handle parameter is a type of function that receive instructions to prepare de variables to be used on topic functions,
// this is necessary because a message need pre processing to handle each event.
func (w *EventHandler) Start(handle func(topics Topics, msg amqp.Delivery), messages <-chan amqp.Delivery) {
	for msg := range messages {
		handle(w.Topics, msg)
	}
}

// AddTopic receive a topic name and the instruction to process the message received for this topic.
func (w *EventHandler) AddTopic(topic string, handleMsgFn func(amqp.Delivery) bool) error {
	if _, ok := w.Topics[topic]; ok {
		return fmt.Errorf("The event %q already exists", topic)
	}

	w.Topics[topic] = handleMsgFn
	return nil
}
