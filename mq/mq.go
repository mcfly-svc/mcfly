package mq

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type MessageChannel interface {
	Send(*amqp.Queue, interface{}) error
	StartDeploy(string, string) error
	CloseConnection() error
	CloseChannel() error
}

func CreateChannel(rabbitMQUrl string) *MsplChannel {

	conn := connect(rabbitMQUrl)
	ch := openChannel(conn)
	buildQueue := declareQueue(ch, "build")

	return &MsplChannel{
		Connection: conn,
		Channel:    ch,
		BuildQueue: buildQueue,
		// other queues...
	}
}

type MsplChannel struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	BuildQueue *amqp.Queue
}

func (rc *MsplChannel) Send(q *amqp.Queue, v interface{}) error {
	json, err := json.Marshal(v)
	if err != nil {
		return NewJsonMarshalError(err, v)
	}
	err = rc.Channel.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/json",
			Body:        json,
		},
	)
	if err != nil {
		return NewQueueSendError(err, json, q.Name)
	}
	return nil
}

func (rc *MsplChannel) Receive(q *amqp.Queue) (<-chan amqp.Delivery, error) {
	msgs, err := rc.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, NewQueueReceiveError(err, q.Name)
	}
	return msgs, nil
}

func (rc *MsplChannel) CloseConnection() error {
	return rc.Connection.Close()
}

func (rc *MsplChannel) CloseChannel() error {
	return rc.Channel.Close()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func connect(rabbitMQUrl string) *amqp.Connection {
	conn, err := amqp.Dial(rabbitMQUrl)
	failOnError(err, "Failed to connect to RabbitMQ")
	return conn
}

func openChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	return ch
}

func declareQueue(ch *amqp.Channel, name string) *amqp.Queue {
	q, err := ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")
	return &q
}
