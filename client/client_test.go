package client_test

import (
	"testing"

	amqp "github.com/cyberlabsai/amqp-client/client"
)

func TestNewClient(t *testing.T) {

	expected := []struct {
		URL         string
		contentType string
	}{
		{
			URL:         "amqp://url1",
			contentType: "text/plain",
		},
		{
			URL:         "amqp://url2",
			contentType: "application/json",
		},
	}

	for i, e := range expected {
		client := amqp.New(e.URL, "text/plain")

		if client.URL != e.URL {
			t.Errorf("Test case %v failed, it should be url equals to [%s] but got [%s]", i, e.URL, client.URL)
		}

		if client.ContentType != e.contentType {
			t.Errorf("Test case %v failed, it should be content type equals to [%s] but got [%s]", i, e.contentType, client.ContentType)
		}
	}
}
