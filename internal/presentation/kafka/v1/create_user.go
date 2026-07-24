package v1

import (
	"context"
	"encoding/json"
	"go-service/internal/core/commands"
	"log"

	kafkago "github.com/segmentio/kafka-go"
)

func (c *Controller) CreateUser(ctx context.Context, msg kafkago.Message) error {
	var cmd commands.CreateUserCommand

	if err := json.Unmarshal(msg.Value, &cmd); err != nil {
		log.Printf(
			"failed to unmarshal create-user message, topic=%s, partition=%d, offset=%d, value=%s, error=%v\n",
			msg.Topic,
			msg.Partition,
			msg.Offset,
			string(msg.Value),
			err,
		)
		return nil
	}

	_, err := c.mediator.HandleCommand(ctx, cmd)
	return err
}
