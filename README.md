<img src="https://raw.githubusercontent.com/Yoonit-Labs/android-yoonit-camera/development/logo_cyberlabs.png" width="300">

# amqp-client

An AMQP client extension to quickly connect and publish messages using a RabbitMQ client.

## Usage

To use the library, connect to an AMQP service passing the variables and will be ready to use.

```go
package main

import amqp "github.com/cyberlabsai/amqp-client/client"

func main() {

    client := amqp.New("my-topic-prefix", "amqp://target_service", "exchange-events", "text/plain")

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