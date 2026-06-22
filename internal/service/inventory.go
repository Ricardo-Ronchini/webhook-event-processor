package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Ricardo-Ronchini/webhook-event-processor/contexts"
)

const (
	inventoryIDMaxLen = 30
	productIDMaxLen   = 30
	skuMaxLen         = 30
	warehouseMaxLen   = 100
)

type Inventory struct {
	InventoryID string    `json:"inventory_id"`
	ProductID   string    `json:"product_id"`
	SKU         string    `json:"sku"`
	Quantity    int       `json:"quantity"`
	Warehouse   string    `json:"warehouse,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (data *Inventory) validation() error {
	if data.InventoryID == "" {
		return fmt.Errorf("id cannot be empty")
	}

	if len(data.InventoryID) > inventoryIDMaxLen {
		return fmt.Errorf("id cannot exceed %d characters", inventoryIDMaxLen)
	}

	if data.ProductID == "" {
		return fmt.Errorf("product_id cannot be empty")
	}

	if len(data.ProductID) > productIDMaxLen {
		return fmt.Errorf("product_id cannot exceed %d characters", productIDMaxLen)
	}

	if data.SKU == "" {
		return fmt.Errorf("sku cannot be empty")
	}

	if len(data.SKU) > skuMaxLen {
		return fmt.Errorf("sku cannot exceed %d characters", skuMaxLen)
	}

	if data.Quantity < 0 {
		return fmt.Errorf("quantity cannot be negative")
	}

	if len(data.Warehouse) > warehouseMaxLen {
		return fmt.Errorf("warehouse cannot exceed %d characters", warehouseMaxLen)
	}

	return nil
}

func (data *Inventory) Get(ctx *contexts.Context) (*Inventory, error) {
	if data.InventoryID == "" {
		return nil, fmt.Errorf("inventory_id cannot be empty")
	}

	query := `
		SELECT 
			id, product_id, sku, quantity, warehouse, created_at, updated_at
		FROM 
			inventory
		WHERE 
			id = $1
	`

	args := []any{
		data.InventoryID,
	}

	row := ctx.App().Database().QueryRow(query, args...)

	var result Inventory

	if err := row.Scan(
		&result.InventoryID,
		&result.ProductID,
		&result.SKU,
		&result.Quantity,
		&result.Warehouse,
		&result.CreatedAt,
		&result.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("unable to continue, scan error")
	}

	return &result, nil
}

func Inventorys(ctx *contexts.Context) (*[]Inventory, error) {
	query := `
		SELECT 
			id, product_id, sku, quantity, warehouse, created_at, updated_at
		FROM 
			inventory
	`

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	rows, err := ctx.App().Database().QueryWithContext(ctxTimeout, query)
	if err != nil {
		return nil, fmt.Errorf("unable to query inventory: %w", err)
	}
	defer rows.Close()

	var result []Inventory

	for rows.Next() {
		var item Inventory

		if err := rows.Scan(
			&item.InventoryID,
			&item.ProductID,
			&item.SKU,
			&item.Quantity,
			&item.Warehouse,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("unable to scan inventory row: %w", err)
		}

		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return &result, nil
}

func (data *Inventory) Create(ctx *contexts.Context) error {
	if err := data.validation(); err != nil {
		return err
	}

	query := `
		INSERT INTO inventory 
			(id, product_id, sku, quantity, warehouse, created_at, updated_at)
		VALUES 
			($1, $2, $3, $4, $5, NOW(), NOW())
	`

	inventoryID := ctx.Helper().GenerateRandomID()

	args := []any{
		inventoryID,
		data.ProductID,
		data.SKU,
		data.Quantity,
		data.Warehouse,
	}

	if _, err := ctx.App().Database().Exec(query, args...); err != nil {
		return fmt.Errorf("unable to create inventory: %w", err)
	}

	return nil
}

func (data *Inventory) Update(ctx *contexts.Context) error {
	if data.InventoryID == "" {
		return fmt.Errorf("inventory_id cannot be empty")
	}

	query := `
		UPDATE inventory
		SET
			product_id  = $2,
			sku         = $3,
			quantity    = $4,
			warehouse   = $5,
			updated_at  = $6
		WHERE 
			id = $1
	`

	args := []any{
		data.InventoryID,
		data.ProductID,
		data.SKU,
		data.Quantity,
		data.Warehouse,
		time.Now(),
	}

	if _, err := ctx.App().Database().Exec(query, args...); err != nil {
		return fmt.Errorf("unable to update inventory: %w", err)
	}

	return nil
}

func DeleteInventory(ctx *contexts.Context, inventoryID string) error {
	if inventoryID == "" {
		return fmt.Errorf("inventory_id cannot be empty")
	}

	query := `
		DELETE FROM inventory
		WHERE inventory_id = $1
	`

	args := []any{
		inventoryID,
	}

	if _, err := ctx.App().Database().Exec(query, args...); err != nil {
		return fmt.Errorf("unable to update inventory: %w", err)
	}

	return nil
}
