package sender

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Sender struct {
	mqconn *amqp.Connection
	ch *amqp.Channel
	q *amqp.Queue
}

func MustNewSender(addr string) *Sender {
	mqconn, err := amqp.Dial(addr)
	if err != nil {
		log.Panicf("Could not connect to RabbitMQ: %v", err)
	}
	ch, err := mqconn.Channel()
	if err != nil {
		log.Panicf("Failed to open a channel %v", err)
	}
	q, err := ch.QueueDeclare(
		"freekill-router", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Panicf("Failed to declare a queue: %v", err)
	}
	return &Sender{
		mqconn: mqconn,
		ch: ch,
		q: &q,
	}
}

func (s *Sender) send(ctx context.Context, body []byte) error {
	err := s.ch.PublishWithContext(
		ctx,
		"",     // exchange
		s.q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing {
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	return err
}

func (s *Sender) Notify(ctx context.Context, command string, data any, connId string) error {
	databody, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// 注意这里的json object只是传消息队列用的，网关会把要传的内容变成数组
	packet := Packet{
		Type: "notify",
		Command: command,
		Data: string(databody),
		ConnectionId: connId,
	}
	body, err := json.Marshal(packet)
	if err != nil {
		return err
	}
	err = s.send(ctx, body)
	return err
}

func (s *Sender) NotifyRaw(ctx context.Context, command string, data string, connId string) error {
	packet := Packet{
		Type: "notify",
		Command: command,
		Data: data,
		ConnectionId: connId,
	}
	body, err := json.Marshal(packet)
	if err != nil {
		return err
	}
	err = s.send(ctx, body)
	return err
}

func (s *Sender) Request(ctx context.Context, command string, data any, connId string, timeout int32, timestamp int64) error {
	databody, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// 注意这里的json object只是传消息队列用的，网关会把要传的内容变成数组
	packet := Packet{
		Type: "request",
		Command: command,
		Data: string(databody),
		ConnectionId: connId,
		Timeout: timeout,
		Timestamp: timestamp,
	}
	body, err := json.Marshal(packet)
	if err != nil {
		return err
	}
	err = s.send(ctx, body)
	return err
}
