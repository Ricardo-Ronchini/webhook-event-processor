package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Ricardo-Ronchini/webhook-event-processor/contexts"
)

const (
	eventTypeMaxLen = 50
)

type InventoryTracks struct {
	TrackID   string          `json:"id"`
	ProductID string          `json:"product_id"`
	EventType string          `json:"event_type"`
	Quantity  int             `json:"quantity"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt time.Time       `json:"created_at"`
}

func (data *InventoryTracks) validation() error {
	if data.TrackID == "" {
		return fmt.Errorf("id cannot be empty")
	}

	if len(data.TrackID) > inventoryIDMaxLen {
		return fmt.Errorf("id cannot exceed %d characters", inventoryIDMaxLen)
	}

	if data.ProductID == "" {
		return fmt.Errorf("product_id cannot be empty")
	}

	if len(data.ProductID) > productIDMaxLen {
		return fmt.Errorf("product_id cannot exceed %d characters", productIDMaxLen)
	}

	if data.EventType == "" {
		return fmt.Errorf("event_type cannot be empty")
	}

	if len(data.EventType) > eventTypeMaxLen {
		return fmt.Errorf("event_type cannot exceed %d characters", eventTypeMaxLen)
	}

	if data.Quantity < 0 {
		return fmt.Errorf("quantity cannot be negative")
	}

	if len(data.Payload) == 0 {
		return fmt.Errorf("payload cannot be empty")
	}

	return nil
}

func (data *InventoryTracks) Get(ctx *contexts.Context) (*InventoryTracks, error) {
	if data.TrackID == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	query := `
		SELECT
			id, product_id, event_type, quantity, payload, created_at
		FROM
			inventory_tracks
		WHERE
			id = $1
	`

	args := []any{
		data.TrackID,
	}

	row := ctx.App().Database().QueryRow(query, args...)

	var result InventoryTracks

	if err := row.Scan(
		&result.TrackID,
		&result.ProductID,
		&result.EventType,
		&result.Quantity,
		&result.Payload,
		&result.CreatedAt,
	); err != nil {
		return nil, fmt.Errorf("unable to scan inventory_tracks row: %w", err)
	}

	return &result, nil
}

func InventoryTrackList(ctx *contexts.Context) (*[]InventoryTracks, error) {
	query := `
		SELECT
			id, product_id, event_type, quantity, payload, created_at
		FROM
			inventory_tracks
	`

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	rows, err := ctx.App().Database().QueryWithContext(ctxTimeout, query)
	if err != nil {
		return nil, fmt.Errorf("unable to query inventory_tracks: %w", err)
	}
	defer rows.Close()

	var result []InventoryTracks

	for rows.Next() {
		var item InventoryTracks

		if err := rows.Scan(
			&item.TrackID,
			&item.ProductID,
			&item.EventType,
			&item.Quantity,
			&item.Payload,
			&item.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("unable to scan inventory_tracks row: %w", err)
		}

		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return &result, nil
}

func (data *InventoryTracks) Create(ctx *contexts.Context) error {
	if err := data.validation(); err != nil {
		return err
	}

	query := `
		INSERT INTO inventory_tracks
			(id, product_id, event_type, quantity, payload, created_at)
		VALUES
			($1, $2, $3, $4, $5, NOW())
	`

	trackID := ctx.Helper().GenerateRandomID()

	args := []any{
		trackID,
		data.ProductID,
		data.EventType,
		data.Quantity,
		data.Payload,
	}

	if _, err := ctx.App().Database().Exec(query, args...); err != nil {
		return fmt.Errorf("unable to create inventory_tracks: %w", err)
	}

	return nil
}
