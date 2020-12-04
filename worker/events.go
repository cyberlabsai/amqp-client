package worker

import (
	"strings"

	"github.com/streadway/amqp"
)

// The library let the developer to create the own Work function,
// to apply rules before the event message processment.

// LastPart is the main example used in our repositories to process common events.
// Ex: onDelete, onCreate, onEdit...
var LastPart = func(events Events, msg amqp.Delivery) {
	parts := strings.Split(msg.RoutingKey, ".")
	lastPart := parts[len(parts)-1]

	if events[lastPart](msg) {
		msg.Ack(false) //must be true?
	}
}
