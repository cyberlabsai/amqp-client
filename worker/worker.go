package worker

import (
	"fmt"

	"github.com/cyberlabsai/amqp-client/client"
	"github.com/streadway/amqp"
)

// Events is a type that contains a list of events and their function.
// The key is the name and the value is what the worker will do when consume a message with this event.
// Ex:
//   event: 	"deleted"
//	 function: 	onDelete(msg amqp.Delivery) bool { // delete from storage }
type Events map[string]func(amqp.Delivery) bool

// Worker consume messages and process according their instructions.
// The service help the worker to connect and bind a topic to receive the messages from network.
// The Events can be added in safe mode with the function AddEvent, or directly.
type Worker struct {
	*client.Service
	Events
}

// New returns a created and connected worker.
// The worker has a service that will declare a new queue
// and link to a topic.
// Before consuming the messages it's necessary to add at least one event
// on your map, after that, the worker can be started successfully.
func New(topicPrefix, address, exchange, contentType, queueName string) (*Worker, error) {

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

	return &Worker{
		Service: s,
	}, nil
}

// AddEvent receive a event name and the instruction to process the message of this event.
func (w *Worker) AddEvent(event string, do func(amqp.Delivery) bool) error {
	if _, ok := w.Events[event]; ok {
		return fmt.Errorf("The event %q already exists", event)
	}

	w.Events[event] = do
	return nil
}

// Work is a type of function that receive instructions to process events.
// This is necessary because somethimes the worker need pre processing in different ways.
type Work func(events Events, msg amqp.Delivery)

// Start initiates the worker to consuming and process the messages.
func (w *Worker) Start(do Work, messages <-chan amqp.Delivery) {
	for msg := range messages {
		// Don't need to wait the process of one message to process another. This is good?
		go do(w.Events, msg)
	}
}
