package client

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

// Service of the amqp to subscribe and publish data.
type Service struct {
	connection  *amqp.Connection
	Channel     *amqp.Channel
	TopicPrefix string
	Address     string
	Exchange    string
	ContentType string
	Mandatory   bool
	Immediate   bool
}

// New returns the client service.
func New(topicPrefix, address, exchange, contentType string) *Service {
	return &Service{
		TopicPrefix: topicPrefix,
		Address:     address,
		Exchange:    exchange,
		ContentType: contentType,
	}
}

// Start the amqp service.
func (s *Service) Start() error {
	connection, err := amqp.Dial(s.Address)
	if err != nil {
		return fmt.Errorf("Error while dialing to AMQP address: %q: %w", s.Address, err)
	}

	s.connection = connection

	s.Channel, err = connection.Channel()
	if err != nil {
		return fmt.Errorf("Error while creating AMQP channel: %w", err)
	}

	return nil
}

// Close do the maintenance to close/clean any connections with the server.
func (s *Service) Close() error {
	defer s.clean()

	if s.Channel != nil {
		err := s.Channel.Close()
		if err != nil {
			return err
		}
	}

	if s.connection != nil {
		return s.connection.Close()
	}

	return nil
}

// Set the connection and channel pointers to empty.
func (s *Service) clean() {
	s.Channel = nil
	s.connection = nil
}

// The restart is called when the service try to publish to a closed channel.
func (s *Service) restart() error {
	// Just empty the connection and start again.
	s.clean()
	return s.Start()
}

// Publish the message to the topic.
func (s *Service) Publish(topic string, body interface{}) error {
	topicName := s.TopicPrefix + "." + topic

	// Marshal into JSON.
	marshalledBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("Error while marshalling JSON data: %w", err)
	}

	// Verify if channel is closed,
	// and will restart the service if it is.
	if s.connection.IsClosed() {
		err = s.restart()
		if err != nil {
			return fmt.Errorf("Error while reconnecting the AMQP service: %w", err)
		}
	}

	return s.Channel.Publish(
		s.Exchange,
		topicName,
		s.Mandatory,
		s.Immediate,
		amqp.Publishing{
			ContentType: s.ContentType,
			Body:        marshalledBody,
		},
	)
}

// Consume messages from a topic.
func (s *Service) Consume(
	consumerName, queueName, topic string) (
	<-chan amqp.Delivery, error) {
	// Declare the queue
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

	// Bind to exchange:
	err = s.Channel.QueueBind(
		queueName,  // queue name
		topic,      // topic
		s.Exchange, // exchange name
		false,      // noWait
		nil,        // Table
	)
	if err != nil {
		return nil, err
	}

	// Open the channel:
	messages, err := s.Channel.Consume(
		queueName,
		consumerName, // consumer name
		false,        // autoAck
		false,        // exclusive
		false,        // noLocal
		false,        // noWait
		nil,          // args
	)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
