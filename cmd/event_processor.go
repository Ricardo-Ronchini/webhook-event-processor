package cmd

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/Ricardo-Ronchini/webhook-event-processor/common"
	"github.com/Ricardo-Ronchini/webhook-event-processor/contexts"
	"github.com/Ricardo-Ronchini/webhook-event-processor/internal/redpanda"
	"github.com/Ricardo-Ronchini/webhook-event-processor/internal/worker"
	"github.com/spf13/cobra"
)

func EventProcessor() *cobra.Command {
	return &cobra.Command{
		Use:   "event-consumer",
		Short: "start event consumer",
		Run: func(cmd *cobra.Command, args []string) {
			c := contexts.NewContext()

			consumer, err := redpanda.NewConsumer(common.ConsumerGroupID, []string{common.MessengerTopicID})
			if err != nil {
				c.App().Logs().Error("[CONSUMER] failed to create consumer: ", err)
				return
			}
			defer consumer.Close()

			// ctx cancelado ao receber SIGINT ou SIGTERM (Ctrl+C ou kill)
			contextBackground, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()

			c.App().Logs().InfoFields("[CONSUMER] started", map[string]any{
				"group_id": common.ConsumerGroupID,
				"topic":    common.MessengerTopicID,
			})

			err = consumer.Poll(contextBackground, func(ctx context.Context, event redpanda.Event) error {
				return worker.ProcessWebhookEvent(c, event)
			})

			if err != nil {
				c.App().Logs().ErrorFields("[CONSUMER] poll error", map[string]any{
					"error": err.Error(),
				})
				return
			}

			c.App().Logs().Info("[CONSUMER] shutdown complete")
		},
	}
}
