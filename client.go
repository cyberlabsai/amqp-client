package main

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

// Service of the amqp to subscribe and publish data.
type Service struct {
	connection  *amqp.Connection
	channel     *amqp.Channel
	topicPrefix string
	address     string
	exchange    string
	mandatory   bool
	immediate   bool
}

// New returns the client service.
func New(topicPrefix, address, exchange string) *Service {
	return &Service{
		topicPrefix: topicPrefix,
		address:     address,
		exchange:    exchange,
	}
}

// Start the amqp service.
func (s *Service) Start() error {
	connection, err := amqp.Dial(s.address)
	if err != nil {
		return fmt.Errorf("Error while dialing to AMQP address: %s", s.address)
	}

	s.connection = connection

	s.channel, err = connection.Channel()
	if err != nil {
		return fmt.Errorf("Error while creating AMQP channel")
	}

	return nil
}

// Close do the maintenance to close/clean any connections with the server.
func (s *Service) Close() error {
	defer s.clean()

	if s.channel != nil {
		err := s.channel.Close()
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
	s.channel = nil
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
	topicName := s.topicPrefix + "." + topic

	// Marshal into JSON.
	marshalledBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("Error while marshalling JSON data")
	}

	// Verify if channel is closed,
	// and will restart the service if it is.
	if s.connection.IsClosed() {
		err = s.restart()
		if err != nil {
			return fmt.Errorf("Error while reconnecting the AMQP service: %s", err.Error())
		}
	}

	return s.channel.Publish(
		s.exchange,
		topicName,
		s.mandatory,
		s.immediate,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        marshalledBody,
		},
	)
}

// This ensures that the service will not be modified once it has been created.

// Exchange return the exchange name.
func (s Service) Exchange() string {
	return s.exchange
}

// TopicPrefix return the topic prefix.
func (s Service) TopicPrefix() string {
	return s.topicPrefix
}

// Address return the amqp address.
func (s Service) Address() string {
	return s.address
}

// These fields can be enabled if necessary.

// EnableMandatory enable the mandatory option.
func (s *Service) EnableMandatory() {
	s.mandatory = true
}

// EnableImmediate enable the immediate option.
func (s *Service) EnableImmediate() {
	s.immediate = true
}
