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
		return fmt.Errorf("Error while dialing to AMQP address: %s", s.Address)
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
	topicName := s.TopicPrefix + "." + topic

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
