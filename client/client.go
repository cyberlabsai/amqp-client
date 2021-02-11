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
	URL         string
	ContentType string
	Mandatory   bool
	Immediate   bool
}

// New returns the client service.
func New(url, contentType string) *Service {
	return &Service{
		URL:         url,
		ContentType: contentType,
	}
}

// Start the amqp service.
func (s *Service) Start() (err error) {
	s.connection, err = amqp.Dial(s.URL)
	if err != nil {
		return fmt.Errorf("Couldn't dialing to amqp url: %q message: %w", s.URL, err)
	}

	s.Channel, err = s.connection.Channel()
	if err != nil {
		return fmt.Errorf("Couldn't open the channel: %w", err)
	}

	return
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
func (s *Service) Publish(exchange, topic string, body interface{}) error {
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
		exchange,
		topic,
		s.Mandatory,
		s.Immediate,
		amqp.Publishing{
			ContentType: s.ContentType,
			Body:        marshalledBody,
		},
	)
}
