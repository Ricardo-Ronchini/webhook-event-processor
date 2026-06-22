package worker

import (
	"fmt"

	"github.com/Ricardo-Ronchini/webhook-event-processor/contexts"
	"github.com/Ricardo-Ronchini/webhook-event-processor/internal/redpanda"
	"github.com/Ricardo-Ronchini/webhook-event-processor/internal/service"
)

func validation(event redpanda.Event) error {
	if event.TenantID == "" {
		return fmt.Errorf("tenant ID is required")
	}

	if event.InventoryID == "" {
		return fmt.Errorf("inventory ID is required")
	}

	if event.ProductID == "" {
		return fmt.Errorf("product ID is required")
	}

	if event.SKU == "" {
		return fmt.Errorf("SKU is required")
	}

	if event.Warehouse == "" {
		return fmt.Errorf("warehouse is required")
	}

	return nil
}

func ProcessWebhookEvent(ctx *contexts.Context, event redpanda.Event) error {
	if err := validation(event); err != nil {
		return err
	}

	// *** responsibilities for handling the received event

	inventory := service.Inventory{
		ProductID: event.ProductID,
		SKU:       event.SKU,
		Quantity:  event.Quantity,
		Warehouse: event.Warehouse,
	}

	if err := inventory.Create(ctx); err != nil {
		return fmt.Errorf("unable to create inventory record: %w", err)
	}

	inventoryTrack := service.InventoryTracks{
		ProductID: event.ProductID,
		EventType: event.Type,
		Quantity:  event.Quantity,
		Payload:   event.Payload(),
	}

	if err := inventoryTrack.Create(ctx); err != nil {
		return fmt.Errorf("unable to create inventory track record: %w", err)
	}

	return nil
}
