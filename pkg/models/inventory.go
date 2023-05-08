package models

type Inventory struct {
	ID         int
	LocationID int
	ItemID     int
	Quantity   float32
	Value      float32
	ValueFIFO  float32
}
