# Go AMQP Client

This is an AMQP client extension to quickly connect and publish messages using a RabbitMQ client.

## Usage

To use the library, connect to an AMQP service passing the variables and will be ready to use.

```go
package main

import amqp "github.com/cyberlabsai/amqp-client"

func main() {

    client := amqp.New("my-topic-prefix", "amqp://target_service", "exchange-events", "text/plain")
    if err := client.Start(); err != nil {
		   return err
    }
	  defer client.Close()
    // The payload can be of any type.
    err := client.Publish("event-done", "payload");
    if err != nil {
        return err
    }
}
```

More datails in unit tests files.

## Testing

```shell
 go test -v ./... 
```

### Contribute with Unit Tests

First, search by uncovered functions.

```shell
  go tool cover -html=coverage.out
```

## Todo list
- [ ] Mock amqp for connection tests?
- [ ] Graceful close the service
- [ ] 100% Unit tests