package models

import "gorm.io/gorm"

type Inventory struct {
	gorm.Model
	ID             uint            `json:"id" gorm:"primaryKey"`
	PlayerID       uint            `json:"player_id"`
	Player         Player          `json:"player" gorm:"foreignKey:PlayerID"`
	InventoryItems []InventoryItem `json:"inventory_items" gorm:"foreignKey:InventoryID"`
}

type InventoryItem struct {
	gorm.Model
	ID          uint      `json:"id" gorm:"primaryKey"`
	InventoryID uint      `json:"inventory_id"`
	Inventory   Inventory `json:"inventory" gorm:"foreignKey:inventory_id"`
	ItemID      uint      `json:"item_id"`
	Quantity    uint      `json:"quantity"`
}
