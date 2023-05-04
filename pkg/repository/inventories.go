package repository

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/models"
)

type (
	inventoryRepository struct {
		conn *sql.DB
	}
)

func NewInventoryRepository(conn *sql.DB) InventoryRepository {
	return &inventoryRepository{conn: conn}
}

func (s *inventoryRepository) GetInventory(tx *sql.Tx, itemID int, locationID int) (m models.Inventory, err error) {
	logger.Debug("get inventory")

	query := `
	SELECT id,  quantity, value, value_fifo
	FROM inventories
	WHERE location_id = ? AND item_id = ?
	`
	m.ItemID = itemID
	m.LocationID = locationID

	err = tx.QueryRow(query, locationID, itemID).Scan(&m.ID, &m.Quantity, &m.Value, &m.ValueFIFO)

	return
}

func (s *inventoryRepository) CreateInventory(tx *sql.Tx, itemID int, locationID int) (m models.Inventory, err error) {
	logger.Debug("create inventory")

	query := `
	INSERT INTO inventories(item_id, location_id)
	VALUES (?, ?)
	`
	m.ItemID = itemID
	m.LocationID = locationID

	var (
		res sql.Result
		id  int64
	)

	res, err = tx.Exec(query, itemID, locationID)
	if err != nil {
		return m, err
	}

	id, err = res.LastInsertId()
	m.ID = int(id)

	return
}

func (s *inventoryRepository) UpdateInventory(tx *sql.Tx, inv models.Inventory) error {
	logger.Debug("update inventory")

	query := `
	UPDATE inventories
	SET item_id = ? , location_id = ?, quantity = ?, value = ?, value_fifo = ?
	WHERE id = ?
	`

	_, err := tx.Exec(query, inv.ItemID, inv.LocationID, inv.Quantity, inv.Value, inv.ValueFIFO, inv.ID)

	return err
}

func (s *inventoryRepository) AddValues(tx *sql.Tx, itemID, locationID int, quantity, val float32) (err error) {
	logger.Debug("inventory add values")

	query := `
	UPDATE inventories
	SET quantity = quantity + ?, value = value + ?, value_fifo = value_fifo + ?
	WHERE item_id = ? AND location_id = ?
	`

	_, err = tx.Exec(query, quantity, val, val, itemID, locationID)

	return
}
