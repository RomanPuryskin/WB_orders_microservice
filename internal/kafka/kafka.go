package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/orders_api/internal/models"
	"github.com/orders_api/internal/service"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	Reader  *kafka.Reader
	Service service.ServiceOrder
	Cfg     *KafkaConfig
}

func NewKafkaConsumer(reader *kafka.Reader, cfg *KafkaConfig, srvc service.ServiceOrder) *KafkaConsumer {
	return &KafkaConsumer{
		Reader:  reader,
		Cfg:     cfg,
		Service: srvc,
	}
}

func NewReader(cfg *KafkaConfig) (*kafka.Reader, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Address, cfg.ExternalPort)

	// попробуем подключиться к kafka
	conn, err := kafka.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("[newReader | failed connect kafka]: %w", err)
	}
	conn.Close()

	readerConfig := kafka.ReaderConfig{
		Brokers: []string{addr},
		Topic:   cfg.Topic,
		GroupID: cfg.Group,
	}

	return kafka.NewReader(readerConfig), nil
}

func (kc *KafkaConsumer) ReadMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := kc.Reader.FetchMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}

				slog.Error("failed to read message from kafka", "error", err)
				continue
			}
			err = kc.ProcessMessage(&msg, ctx)
			if err != nil {
				slog.Error("Failed to ProcessMessage", "error", err)
				// закоммитим невалидное сообщение, чтобы больше его не читать
				err = kc.Reader.CommitMessages(ctx, msg)
				if err != nil {
					slog.Error("Failed to Commit Message", "error", err)
				}
				continue
			}

			// закоммитим сообщение (оно удачно загрузилось в БД)
			err = kc.Reader.CommitMessages(ctx, msg)
			if err != nil {
				slog.Error("Failed to Commit Message", "error", err)
				continue
			}
		}
	}
}

func (kc *KafkaConsumer) ProcessMessage(msg *kafka.Message, ctx context.Context) error {
	var newOrder models.Order
	err := json.Unmarshal(msg.Value, &newOrder)
	if err != nil {
		return fmt.Errorf("[ProcessMessage| failed to Unmarshal]: %w", err)
	}

	_, err = kc.Service.SetOrder(&newOrder)
	if err != nil {
		return fmt.Errorf("[ProcessMessage| failed to SetOrder]: %w", err)
	}

	slog.Info("Success readed msg and set order", "order_uid", newOrder.OrderUID)
	return nil
}
