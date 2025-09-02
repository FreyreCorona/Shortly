// Package rabbitmq used to consume queued messages and cache it
package rabbitmq

import (
	"context"
	"encoding/json"
	"log"

	"github.com/FreyreCorona/Shortly/src/redirect_svc/internal/application"
	"github.com/FreyreCorona/Shortly/src/redirect_svc/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	service application.SetURLService
	ch      *amqp.Channel
	q       *amqp.Queue
}

func NewConsumer(service application.SetURLService, address string) (*Consumer, error) {
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
	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		return nil, err
	}
	err = ch.QueueBind(q.Name, "created", "cache", false, nil)
	if err != nil {
		return nil, err
	}

	return &Consumer{service: service, ch: ch, q: &q}, nil
}

func (c *Consumer) Listen(ctx context.Context) error {
	msgs, err := c.ch.Consume(c.q.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case d, ok := <-msgs:
				if !ok {
					log.Println("Channel closed")
					return
				}

				var url struct {
					RawURL    string `json:"raw_url"`
					ShortCode string `json:"short_code"`
				}

				result := d.Body
				err = json.Unmarshal(result, &url)
				if err != nil {
					log.Printf("error unmarshalling the result :%v", err)
					d.Nack(false, false)
				}

				u := domain.URL{RawURL: url.RawURL, ShortCode: url.ShortCode}
				err = c.service.SetURL(u)
				if err != nil {
					log.Printf("error caching the url :%v", err)
				}

				d.Ack(false)
			case <-ctx.Done():
				log.Println("Stoping consumer")
				return
			}
		}
	}()
	<-ctx.Done()
	if err := c.ch.Close(); err != nil {
		log.Printf("error closing channel: %v", err)
		return err
	}

	return nil
}

func (c *Consumer) SetURL(url domain.URL) error {
	return c.service.SetURL(url)
}
