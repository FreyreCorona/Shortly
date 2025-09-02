// Package rabbitmq used to enqueue messages
package rabbitmq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/FreyreCorona/Shortly/src/shortener_svc/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ProducerPublisher struct {
	ch *amqp.Channel
}

func NewProducerPublisher(address string) (*ProducerPublisher, error) {
	conn, err := amqp.Dial(address)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	err = ch.ExchangeDeclare("cache", "direct", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &ProducerPublisher{ch: ch}, nil
}

func (p *ProducerPublisher) Publish(url domain.URL) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var u struct {
		RawURL    string `json:"raw_url"`
		ShortCode string `json:"short_code"`
	}
	u.RawURL = url.RawURL
	u.ShortCode = url.ShortCode

	result, err := json.Marshal(u)
	if err != nil {
		return err
	}

	err = p.ch.PublishWithContext(ctx, "cache", "created", false, false, amqp.Publishing{ContentType: "application/json", Body: result})
	if err != nil {
		return err
	}

	return nil
}
