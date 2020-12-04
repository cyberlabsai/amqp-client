package client_test

import (
	"testing"

	amqp "github.com/cyberlabsai/amqp-client/client"
)

func TestNewClient(t *testing.T) {

	expected := []struct {
		topicPrefix string
		address     string
		exchange    string
	}{
		{
			topicPrefix: "test-topic",
			address:     "amqp://address",
			exchange:    "test-exchange",
		},
		{
			topicPrefix: "test-topic2",
			address:     "amqp://address2",
			exchange:    "",
		},
	}

	for i, e := range expected {
		client := amqp.New(e.topicPrefix, e.address, e.exchange, "text/plain")

		if client.TopicPrefix != e.topicPrefix {
			t.Errorf("Test case %v failed, it should be topic prefix equals to [%s] but got [%s]", i, e.topicPrefix, client.TopicPrefix)
		}

		if client.Address != e.address {
			t.Errorf("Test case %v failed, it should be address equals to [%s] but got [%s]", i, e.exchange, client.Address)
		}

		if client.Exchange != e.exchange {
			t.Errorf("Test case %v failed, it should be exchange [%s] but got [%s]", i, e.exchange, client.Exchange)
		}
	}
}
